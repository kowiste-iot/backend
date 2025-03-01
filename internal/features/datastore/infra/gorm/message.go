package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/features/datastore/domain"
	"backend/internal/features/datastore/domain/command"
	ingestDomain "backend/internal/features/ingest/domain"
	baseCmd "backend/shared/base/command"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MessageModel represents the database model for stored messages
type MessageModel struct {
	ID        string         `gorm:"primaryKey;type:uuid"`
	MeasureID string         `gorm:"uniqueIndex:idx_tenant_branch_measure_time;not null"`
	TenantID  string         `gorm:"uniqueIndex:idx_tenant_branch_measure_time;not null"`
	BranchID  string         `gorm:"uniqueIndex:idx_tenant_branch_measure_time;not null"`
	Time      time.Time      `gorm:"uniqueIndex:idx_tenant_branch_measure_time;not null"`
	Data      datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedAt time.Time      `gorm:"not null"`
}

func (MessageModel) TableName() string {
	return "message_store"
}
// BeforeCreate will set a UUID rather than numeric ID
func (m *MessageModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		uuid, err := uuid.NewV7()
		if err != nil {
			return err
		}
		m.ID = uuid.String()
	}	
	return nil
}
func (m MessageModel) DataToInterface() (result map[string]any, err error) {

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

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.MessageRepository {
	db.AutoMigrate(MessageModel{})

	return &repository{
		db: db,
	}
}

func (r *repository) Store(ctx context.Context, msg *ingestDomain.Message) error {
	b, err := msg.DataToBytes()
	if err != nil {
		return fmt.Errorf("failed to marshal message data: %w", err)
	}
	model := &MessageModel{
		MeasureID: msg.ID,
		TenantID:  msg.TenantID,
		BranchID:  msg.BranchID,
		Time:      msg.Time,
		Data:      datatypes.JSON(b),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to store message: %w", err)
	}

	return nil
}

func (r *repository) StoreBatch(ctx context.Context, msgs []*ingestDomain.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	models := make([]*MessageModel, len(msgs))

	for i, msg := range msgs {
		b, err := msg.DataToBytes()
		if err != nil {
			return fmt.Errorf("failed to marshal message data: %w", err)
		}
		models[i] = &MessageModel{
			MeasureID:        msg.ID,
			TenantID:  msg.TenantID,
			BranchID:  msg.BranchID,
			Time:      msg.Time,
			Data:      datatypes.JSON(b),
		}
	}

	// Use CreateInBatches for better performance with large batches
	if err := r.db.WithContext(ctx).CreateInBatches(models, len(models)).Error; err != nil {
		return fmt.Errorf("failed to store message batch: %w", err)
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, input *baseCmd.BaseInput, id string) (*ingestDomain.Message, error) {
	var model MessageModel

	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND branch_id = ?", id, input.TenantDomain, input.BranchName).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMessageNotFound
		}
		return nil, fmt.Errorf("failed to find message: %w", err)
	}

	return modelToDomain(&model)
}

func (r *repository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*ingestDomain.Message, error) {
	var models []MessageModel

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ?", input.TenantDomain, input.BranchName).
		Order("time DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	return modelsToDomain(models)
}

func (r *repository) FindByTimeRange(ctx context.Context, input *command.TimeRangeInput) ([]*ingestDomain.Message, error) {
	var models []MessageModel

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ? AND time BETWEEN ? AND ?",
			input.TenantDomain,
			input.BranchName,
			input.StartTime,
			input.EndTime,
		).
		Order("time DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find messages in time range: %w", err)
	}

	return modelsToDomain(models)
}

// Helper functions to convert between models and domain objects
func modelToDomain(model *MessageModel) (*ingestDomain.Message, error) {
	value, err := model.DataToInterface()
	if err != nil {
		return nil, err
	}
	msg := &ingestDomain.Message{
		ID:       model.ID,
		TenantID: model.TenantID,
		BranchID: model.BranchID,
		Time:     model.Time,
		Data:     value,
	}

	return msg, nil
}

func modelsToDomain(models []MessageModel) ([]*ingestDomain.Message, error) {
	stores := make([]*ingestDomain.Message, len(models))
	for i, model := range models {
		store, err := modelToDomain(&model)
		if err != nil {
			return nil, err
		}
		stores[i] = store
	}
	return stores, nil
}
