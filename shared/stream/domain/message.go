package domain

import (
	"backend/internal/features/ingest/domain"
	"time"

	"github.com/google/uuid"
)

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

func NewFromIngest(input *domain.Message) (m *Message, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	m = &Message{
		ID:        id.String(),
		TenantID:  input.TenantID,
		Topic:     domain.TopicIngest,
		Data:      input,
		Timestamp: input.Time,
		Event:     domain.EventIngest,
	}
	if m.Timestamp.IsZero() {
		m.Timestamp = time.Now()
	}
	return
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
