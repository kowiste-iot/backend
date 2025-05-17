package repository

import (
	ingestDomain "backend/internal/features/ingest/domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MessageModelDB struct {
	ID        string         `gorm:"primaryKey;type:uuid"`
	MeasureID string         `gorm:"uniqueIndex:idx_tenant_branch_measure_time;not null"`
	Time      time.Time      `gorm:"uniqueIndex:idx_tenant_branch_measure_time;not null"`
	Data      datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedAt time.Time      `gorm:"not null"`
}

func (MessageModelDB) TableName() string {
	return "message_store"
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *MessageModelDB) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		uuid, err := uuid.NewV7()
		if err != nil {
			return err
		}
		m.ID = uuid.String()
	}
	return nil
}
func (m MessageModelDB) DataToInterface() (result map[string]any, err error) {

	// If Data is empty, return an empty map
	if len(m.Data) == 0 {
		return make(map[string]interface{}), nil
	}

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(m.Data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	return result, nil
}

// Mapper domain model to db model
func messageMapper(data *ingestDomain.Message) (dbMessage *MessageModelDB, err error) {
	b, err := data.DataToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message data: %w", err)
	}
	dbMessage = &MessageModelDB{
		MeasureID: data.ID,
		Time:      data.Time,
		Data:      datatypes.JSON(b),
	}
	return
}
