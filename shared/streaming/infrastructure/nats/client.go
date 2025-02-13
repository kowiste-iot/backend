package nats

import (
	"backend/shared/streaming/domain"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	maxRetries = 5
	baseDelay  = 100 * time.Millisecond
	maxDelay   = 5 * time.Second
)

type subscription struct {
	sub     *nats.Subscription
	handler domain.MessageHandler
}

type Client struct {
	conn          *Connection
	subscriptions map[string]subscription
	mu            sync.RWMutex
}

func NewClient(conn *Connection) *Client {
	return &Client{
		conn:          conn,
		subscriptions: make(map[string]subscription),
	}
}

func (c *Client) Subscribe(ctx context.Context, subjectGen domain.SubjectGenerator, handler domain.MessageHandler, params ...string) error {
	subject := subjectGen.Generate(params...)

	sub, err := c.conn.GetConn().Subscribe(subject, c.createMessageHandler(handler))
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	c.mu.Lock()
	c.subscriptions[subject] = subscription{
		sub:     sub,
		handler: handler,
	}
	c.mu.Unlock()

	return nil
}

func (c *Client) createMessageHandler(handler domain.MessageHandler) nats.MsgHandler {
	return func(msg *nats.Msg) {
		var notification domain.Message
		if err := json.Unmarshal(msg.Data, &notification); err != nil {
			return
		}
		handler.Handle(notification)
	}
}

func (c *Client) resubscribeWithBackoff(subject string) error {
	var err error

	c.mu.RLock()
	subs, exists := c.subscriptions[subject]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("subscription not found for subject: %s", subject)
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		delay := time.Duration(math.Min(
			float64(baseDelay)*math.Pow(2, float64(attempt)),
			float64(maxDelay),
		))
		time.Sleep(delay)

		sub, err := c.conn.GetConn().Subscribe(subject, c.createMessageHandler(subs.handler))
		if err == nil {
			c.mu.Lock()
			c.subscriptions[subject] = subscription{
				sub:     sub,
				handler: subs.handler,
			}
			c.mu.Unlock()
			return nil
		}
	}

	return fmt.Errorf("failed to resubscribe after %d attempts: %w", maxRetries, err)
}
func (c *Client) Unsubscribe(subjectGen domain.SubjectGenerator, params ...string) error {
	subject := subjectGen.Generate(params...)

	c.mu.Lock()
	defer c.mu.Unlock()

	if sub, ok := c.subscriptions[subject]; ok {
		if err := sub.sub.Unsubscribe(); err != nil {
			return err
		}
		delete(c.subscriptions, subject)
	}

	return nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, subscription := range c.subscriptions {
		if err := subscription.sub.Unsubscribe(); err != nil {
			continue
		}
	}

	c.subscriptions = make(map[string]subscription)
	return nil
}

func (c *Client) Publish(ctx context.Context, subjectGen domain.SubjectGenerator, msg domain.Message, params ...string) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	subject := subjectGen.Generate(params...)
	return c.conn.GetConn().Publish(subject, data)
}

func (c *Client) HandleReconnect(ctx context.Context) {
	c.conn.GetConn().SetReconnectHandler(func(nc *nats.Conn) {
		c.mu.RLock()
		subjects := make([]string, 0, len(c.subscriptions))
		for subject := range c.subscriptions {
			subjects = append(subjects, subject)
		}
		c.mu.RUnlock()

		for _, subject := range subjects {
			if err := c.resubscribeWithBackoff(subject); err != nil {
				continue
			}
		}
	})
}
