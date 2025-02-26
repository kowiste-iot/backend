package domain

import (
	"backend/shared/validator"
	"encoding/json"
	"time"
)

type Message struct {
	ID       string                 `json:"id" validate:"uuidv7"`
	TenantID string                 `json:"tenant" validate:"required,min=3,max=255"`
	BranchID string                 `json:"branch" validate:"required,min=3,max=255"`
	Time     time.Time              `json:"time"`
	Data     map[string]interface{} `json:"data" validate:"required"`
}

// Validate checks if the message contains all required fields
func (m *Message) Validate() error {
	return validator.Validate(m)
}

// ToBytes serializes the message into bytes
func (m *Message) ToBytes() ([]byte, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
func (m *Message) DataToBytes() ([]byte, error) {

	return json.Marshal(m.Data)
}
