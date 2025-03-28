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

//mqtt
// {
//   "id":"019545fa-101a-7d62-8197-fe6930b5f058",
//   "tenant":"019545fa-6ae7-7d61-bd7d-a7dc67a5403e",
//   "branch":"019545fa-b1a7-7845-80cc-b0b55a9a4d42",
//   "data":{
//     "hello":3
//   }
// }

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
