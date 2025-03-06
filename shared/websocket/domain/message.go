// shared/websocket/domain/message.go
package domain

import "time"

// MessageType defines the type of websocket message
type MessageType string

const (
	TypeSubscribe MessageType = "subscribe"
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

// Alert specific message content
type AlertContent struct {
	Source   string      `json:"source"`
	Severity string      `json:"severity"`
	Details  interface{} `json:"details"`
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
