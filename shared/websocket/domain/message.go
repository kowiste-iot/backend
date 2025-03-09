// shared/websocket/domain/message.go
package domain

import (
	streamDomain "backend/shared/stream/domain"
	"context"
	"time"
)

// MessageType defines the type of websocket message
type MessageType string

const (
	TypeSubscribe   MessageType = "subscribe"
	TypeUnsubscribe MessageType = "unsubscribe"
	TypeGetValue    MessageType = "getValue"
)

// Message represents a websocket message
type Message struct {
	ID        string      `json:"id"`
	TenantID  string      `json:"tenantID"`
	BranchID  string      `json:"branchID"`
	UserID    string      `json:"userID"`
	Type      MessageType `json:"type"`
	Content   any         `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
}

// NewMessage creates a new message with default values
func NewMessage(tenantID, userID string, msgType MessageType) *Message {
	return &Message{
		TenantID:  tenantID,
		UserID:    userID,
		Type:      msgType,
		CreatedAt: time.Now(),
	}
}

// MessageHandler defines a function that handles a specific message type
type MessageHandler func(ctx context.Context, message *streamDomain.WireMessage) error
