// shared/websocket/domain/message.go
package domain

import "time"

// MessageType defines the type of websocket message
type MessageType string

const (
	TypeAlert  MessageType = "alert"
	TypeSystem MessageType = "system"
	TypeUpdate MessageType = "update"
)

// MessagePriority defines the priority level of the message
type MessagePriority string

const (
	PriorityLow    MessagePriority = "low"
	PriorityMedium MessagePriority = "medium"
	PriorityHigh   MessagePriority = "high"
)

// Message represents a websocket message
type Message struct {
	ID        string          `json:"id"`
	TenantID  string          `json:"tenant_id"`
	UserID    string          `json:"user_id"`
	Type      MessageType     `json:"type"`
	Priority  MessagePriority `json:"priority"`
	Title     string          `json:"title"`
	Content   interface{}     `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	ReadAt    *time.Time      `json:"read_at,omitempty"`
}

// Alert specific message content
type AlertContent struct {
	Source   string      `json:"source"`
	Severity string      `json:"severity"`
	Details  interface{} `json:"details"`
}

// System specific message content
type SystemContent struct {
	Category string      `json:"category"`
	Action   string      `json:"action"`
	Details  interface{} `json:"details"`
}

// Update specific message content
type UpdateContent struct {
	Entity  string      `json:"entity"`
	Action  string      `json:"action"`
	Details interface{} `json:"details"`
}

// NewMessage creates a new message with default values
func NewMessage(tenantID, userID string, msgType MessageType, priority MessagePriority) *Message {
	return &Message{
		TenantID:  tenantID,
		UserID:    userID,
		Type:      msgType,
		Priority:  priority,
		CreatedAt: time.Now(),
	}
}
