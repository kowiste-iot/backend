package app

import (
	"backend/shared/websocket/domain"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client struct update
type Client struct {
	hub             *Hub
	conn            *websocket.Conn
	send            chan []byte
	tenantID        string
	userID          string
	mu              sync.Mutex
	commandHandlers map[domain.MessageType]CommandHandler
}

// Define the command handler type
type CommandHandler func(message *domain.Message) error

// NewClient function update
func NewClient(hub *Hub, conn *websocket.Conn, tenantID, userID string) *Client {
	client := &Client{
		hub:             hub,
		conn:            conn,
		send:            make(chan []byte, 256),
		tenantID:        tenantID,
		userID:          userID,
		commandHandlers: make(map[domain.MessageType]CommandHandler),
	}

	// Register default command handlers
	client.registerDefaultHandlers()

	return client
}

// Register default command handlers
func (c *Client) registerDefaultHandlers() {
	c.commandHandlers[domain.TypeSubscribe] = c.handleSubscribe
	c.commandHandlers[domain.TypeUnsubscribe] = c.handleUnsubscribe
	c.commandHandlers[domain.TypeGetValue] = c.handleGetCurrentValue

}

// ReadPump with updated command handling
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error read: %v\n", err)
			}
			break
		}

		// Process the message based on its type
		var clientMsg *domain.Message
		if err := json.Unmarshal(data, &clientMsg); err != nil {
			fmt.Printf("error unmarshaling message: %v\n", err)
			continue
		}

		if handler, exists := c.commandHandlers[clientMsg.Type]; exists {
			if err := handler(clientMsg); err != nil {
				fmt.Printf("error handling command: %v\n", err)
			}
			continue
		}
		fmt.Printf("unknown message type: %s\n", clientMsg.Type)

	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.mu.Lock()
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()
			if err != nil {
				return
			}
		case <-ticker.C:
			c.mu.Lock()
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			c.mu.Unlock()
			if err != nil {
				return
			}
		}
	}
}
