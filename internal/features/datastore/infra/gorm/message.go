package repository

import (
	"context"
	"fmt"

	"backend/internal/features/datastore/domain"
	"backend/internal/features/datastore/domain/command"
	ingestDomain "backend/internal/features/ingest/domain"
	gormhelper "backend/shared/gorm"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.MessageRepository {

	return &repository{
		db: db,
	}
}

func (r *repository) Store(ctx context.Context, msg *ingestDomain.Message) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	model, err := messageMapper(msg)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to store message: %w", err)
	}

	return nil
}

func (r *repository) StoreBatch(ctx context.Context, msgs []*ingestDomain.Message) (err error) {
	if len(msgs) == 0 {
		return nil
	}

	models := make([]*MessageModelDB, len(msgs))

	for i, msg := range msgs {
		model, err := messageMapper(msg)
		if err != nil {
			return err
		}
		models[i] = model
	}
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	// Use CreateInBatches for better performance with large batches
	if err := db.WithContext(ctx).CreateInBatches(models, len(models)).Error; err != nil {
		return fmt.Errorf("failed to store message batch: %w", err)
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*ingestDomain.Message, error) {
	var model MessageModelDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).
		Where("id = ? ", id).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMessageNotFound
		}
		return nil, fmt.Errorf("failed to find message: %w", err)
	}

	return modelToDomain(&model)
}

func (r *repository) FindAll(ctx context.Context) ([]*ingestDomain.Message, error) {
	var models []MessageModelDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).
		Order("time DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	return modelsToDomain(models)
}

func (r *repository) FindByTimeRange(ctx context.Context, input *command.TimeRangeInput) ([]*ingestDomain.Message, error) {
	var models []MessageModelDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).
		Where("time BETWEEN ? AND ?",
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
func modelToDomain(model *MessageModelDB) (*ingestDomain.Message, error) {
	value, err := model.DataToInterface()
	if err != nil {
		return nil, err
	}
	msg := &ingestDomain.Message{
		ID:       model.ID,
		TenantID: "model.TenantID",
		BranchID: "model.BranchID",
		Time:     model.Time,
		Data:     value,
	}

	return msg, nil
}

func modelsToDomain(models []MessageModelDB) ([]*ingestDomain.Message, error) {
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
