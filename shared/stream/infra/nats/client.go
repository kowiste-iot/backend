package nats

import (
	"backend/shared/stream/domain"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn         *Connection
	config       *domain.StreamConfig
	subscribers  map[string]*nats.Subscription
	subscriberMu sync.RWMutex
	repository   domain.MessageRepository
}

func NewClient(conn *Connection, config *domain.StreamConfig, repo domain.MessageRepository) *NatsClient {
	return &NatsClient{
		conn:        conn,
		config:      config,
		subscribers: make(map[string]*nats.Subscription),
		repository:  repo,
	}
}

func (c *NatsClient) Connect() error {
	return c.conn.Connect()
}

func (c *NatsClient) Publish(ctx context.Context, msg *domain.Message) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to NATS")
	}

	// Convert MessageData to bytes
	data, err := msg.Data.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to convert message data to bytes: %w", err)
	}

	// Create wire format
	wireMsg := domain.WireMessage{
		ID:        msg.ID,
		Topic:     msg.Topic,
		Data:      data,
		Timestamp: msg.Timestamp,
		Event:     msg.Event,
	}

	// If persistence enabled, save original message
	if c.config.PersistMessage && c.repository != nil {
		msg.Status = domain.MessageStatusPending
		if err := c.repository.Save(msg); err != nil {
			return fmt.Errorf("failed to persist message: %w", err)
		}
	}

	// Marshal wire message
	wireData, err := json.Marshal(wireMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal wire message: %w", err)
	}

	// Publish to NATS
	err = c.conn.GetConn().Publish(msg.Topic, wireData)
	if err != nil {
		if c.repository != nil {
			c.repository.UpdateStatus(msg.ID, domain.MessageStatusFailed)
		}
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Update status if persistence is enabled
	if c.config.PersistMessage && c.repository != nil {
		if err := c.repository.UpdateStatus(msg.ID, domain.MessageStatusSent); err != nil {
			return fmt.Errorf("failed to update message status: %w", err)
		}
	}

	return nil
}

func (c *NatsClient) Subscribe(ctx context.Context, topic string, handler domain.MessageHandler) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to NATS")
	}

	c.subscriberMu.Lock()
	defer c.subscriberMu.Unlock()

	if _, exists := c.subscribers[topic]; exists {
		return fmt.Errorf("already subscribed to topic: %s", topic)
	}

	sub, err := c.conn.GetConn().Subscribe(topic, func(m *nats.Msg) {
		var wireMsg domain.WireMessage
		if err := json.Unmarshal(m.Data, &wireMsg); err != nil {
			// Handle error, maybe through a error channel or logger
			return
		}
		
		if err := handler(ctx, &wireMsg); err != nil {
			// Handle error
			return
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}
	fmt.Println("subscribed to topic: ", topic)
	c.subscribers[topic] = sub
	return nil
}

func (c *NatsClient) Unsubscribe(ctx context.Context, topic string) error {
	c.subscriberMu.Lock()
	defer c.subscriberMu.Unlock()

	sub, exists := c.subscribers[topic]
	if !exists {
		return fmt.Errorf("no subscription found for topic: %s", topic)
	}

	if err := sub.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe from topic: %w", err)
	}

	delete(c.subscribers, topic)
	return nil
}

func (c *NatsClient) Close() error {
	c.subscriberMu.Lock()
	defer c.subscriberMu.Unlock()

	for topic, sub := range c.subscribers {
		if err := sub.Unsubscribe(); err != nil {
			return fmt.Errorf("failed to unsubscribe from topic %s: %w", topic, err)
		}
		delete(c.subscribers, topic)
	}

	return c.conn.Close()
}

func (c *NatsClient) IsConnected() bool {
	return c.conn.IsConnected()
}

func (c *NatsClient) GetConfig() *domain.StreamConfig {
	return c.config
}
