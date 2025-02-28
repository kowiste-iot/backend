package domain

import (
	"backend/internal/features/ingest/domain"
	"context"
	"encoding/json"
	"reflect"
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

func (m Message) DataToModel(model interface{}) error {
	// First, convert the Data to bytes
	dataBytes, err := m.Data.ToBytes()
	if err != nil {
		return err
	}
	// Check if model is a pointer
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		return json.Unmarshal(dataBytes, &model)
	}
	// Unmarshal the bytes into the provided model
	return json.Unmarshal(dataBytes, model)
}

func NewFromIngest(input *domain.Message) (m *Message, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	if input.Time.IsZero() {
		input.Time = time.Now()
	}
	m = &Message{
		ID:        id.String(),
		TenantID:  input.TenantID,
		Topic:     domain.TopicIngest,
		Data:      input,
		Timestamp: time.Now(),
		Event:     domain.EventIngest,
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

func (m WireMessage) DataToModel(model interface{}) error {
	// Check if model is a pointer
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		return json.Unmarshal(m.Data, &model)
	}
	// Unmarshal the bytes into the provided model
	return json.Unmarshal(m.Data, model)
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

type MessageHandler func(ctx context.Context, msg *WireMessage) error

type MessageRepository interface {
	Save(message *Message) error
	UpdateStatus(id string, status MessageStatus) error
	GetPendingMessages() ([]*StoredMessage, error)
}
