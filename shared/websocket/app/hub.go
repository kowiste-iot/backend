// shared/websocket/app/hub.go
package app

import (
	"context"
	"sync"
)

type Hub struct {
	// Clients map with tenant isolation
	// map[tenantID]map[userID]*Client
	clients    map[string]map[string]*Client
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.tenantID]; !ok {
				h.clients[client.tenantID] = make(map[string]*Client)
			}
			h.clients[client.tenantID][client.userID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.tenantID]; ok {
				if _, ok := h.clients[client.tenantID][client.userID]; ok {
					delete(h.clients[client.tenantID], client.userID)
					close(client.send)
					// If tenant has no more users, clean up
					if len(h.clients[client.tenantID]) == 0 {
						delete(h.clients, client.tenantID)
					}
				}
			}
			h.mu.Unlock()

		case <-ctx.Done():
			return
		}
	}
}

// SendToUser sends a message to a specific user in a tenant
func (h *Hub) SendToUser(tenantID, userID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if tenantClients, ok := h.clients[tenantID]; ok {
		if client, ok := tenantClients[userID]; ok {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(tenantClients, userID)
			}
		}
	}
}

// GetConnectedUsers returns all connected users for a tenant
func (h *Hub) GetConnectedUsers(tenantID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var users []string
	if tenantClients, ok := h.clients[tenantID]; ok {
		for userID := range tenantClients {
			users = append(users, userID)
		}
	}
	return users
}

// IsUserConnected checks if a specific user is connected
func (h *Hub) IsUserConnected(tenantID, userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if tenantClients, ok := h.clients[tenantID]; ok {
		_, exists := tenantClients[userID]
		return exists
	}
	return false
}

func (h *Hub) RegisterClient(client *Client) {
    h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
    h.unregister <- client
}