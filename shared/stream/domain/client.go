package domain

import "context"

type StreamClient interface {
	// Connect establishes the connection with the streaming server
	Connect() error

	// Publish sends a message to the specified topic
	Publish(ctx context.Context, msg *Message) error

	// Subscribe registers a handler for a specific topic
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error

	// Unsubscribe removes the subscription for a topic
	Unsubscribe(ctx context.Context, topic string) error

	// Close closes the connection and cleans up resources
	Close() error

	// IsConnected returns the current connection status
	IsConnected() bool

	// GetConfig returns the current client configuration
	GetConfig() *StreamConfig
}

// For dependency injection in services
type StreamClientFactory interface {
	CreateClient(config *StreamConfig) (StreamClient, error)
}
