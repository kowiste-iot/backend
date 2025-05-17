package repository

import (
	"backend/internal/features/action/domain"
	gormhelper "backend/shared/gorm"
	"backend/shared/pagination"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type actionRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.ActionRepository {
	return &actionRepository{db: db}
}

func (r *actionRepository) Create(ctx context.Context, input *domain.Action) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Create(actiondMapper(input)).Error
}

func (r *actionRepository) Update(ctx context.Context, input *domain.Action) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Updates(actiondMapper(input)).Error
}

func (r *actionRepository) FindByID(ctx context.Context, actionID string) (*domain.Action, error) {
	var dbAction ActionDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", actionID).
		First(&dbAction).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbAction.ID,
		"dbAction.TenantID",
		"dbAction.BranchID",
		dbAction.Name,
		dbAction.Description,
		dbAction.Parent,
		dbAction.Enabled,
		dbAction.UpdatedAt,
		&dbAction.DeletedAt.Time,
	), nil

}

func (r *actionRepository) FindAll(ctx context.Context) ([]*domain.Action, error) {
	var dbActions []ActionDB

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	var total int64
	err = db.Model(&ActionDB{}).Where(gormhelper.DeleteFilter()).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbActions).Error
	if err != nil {
		return nil, err
	}

	actions := make([]*domain.Action, len(dbActions))
	for i, dbAction := range dbActions {
		actions[i] = domain.NewFromRepository(
			dbAction.ID,
			"dbAction.TenantID",
			"dbAction.BranchID",
			dbAction.Name,
			dbAction.Description,
			dbAction.Parent,
			dbAction.Enabled,
			dbAction.UpdatedAt,
			&dbAction.DeletedAt.Time,
		)
	}
	return actions, nil
}

func (r *actionRepository) Remove(ctx context.Context, actionID string) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	resp := db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", actionID).Delete(&ActionDB{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}
