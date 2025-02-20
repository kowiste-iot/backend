package domain

import "time"

type MessageData interface {
	Validate() error
	ToBytes() ([]byte, error)
}

// Message is used by services to send data
type Message struct {
	ID        string
	TenantID  string
	UserID    string
	Topic     string
	Data      MessageData
	Timestamp time.Time
	Event     string
	Status    MessageStatus
}

// WireMessage is the format used for NATS transmission
type WireMessage struct {
	ID        string
	Topic     string
	Data      []byte
	Timestamp time.Time
	Event     string
}

// StoredMessage represents a message as stored in the database
type StoredMessage struct {
	ID        string
	Topic     string
	Data      []byte
	Timestamp time.Time
	Event     string
	Status    MessageStatus
}

type MessageStatus string

const (
	MessageStatusPending MessageStatus = "pending"
	MessageStatusSent    MessageStatus = "sent"
	MessageStatusFailed  MessageStatus = "failed"
)

type MessageHandler func(msg *WireMessage) error

type MessageRepository interface {
	Save(message *Message) error
	UpdateStatus(id string, status MessageStatus) error
	GetPendingMessages() ([]*StoredMessage, error)
}
