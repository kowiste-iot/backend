package repository

import (
	"context"
	"ddd/internal/features/action/domain"
	baseCmd "ddd/shared/base/command"
	gormhelper "ddd/shared/gorm"
	"ddd/shared/pagination"
	"errors"
	"time"

	"gorm.io/gorm"
)

type actionRepository struct {
	db *gorm.DB
}

type Action struct {
	ID          string `gorm:"primaryKey"`
	TenantID    string `gorm:"index"`
	BranchID    string `gorm:"index"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Enabled     bool
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func NewRepository(db *gorm.DB) domain.ActionRepository {
	db.AutoMigrate(&Action{})
	return &actionRepository{db: db}
}

func (r *actionRepository) Create(ctx context.Context, input *domain.Action) error {
	dbAction := Action{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Create(&dbAction).Error
}

func (r *actionRepository) Update(ctx context.Context, input *domain.Action) error {
	dbAction := Action{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Updates(&dbAction).Error
}

func (r *actionRepository) FindByID(ctx context.Context, input *baseCmd.BaseInput, actionID string) (*domain.Action, error) {
	var dbAction Action

	err := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", actionID).
		First(&dbAction).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbAction.ID,
		dbAction.TenantID,
		dbAction.BranchID,
		dbAction.Name,
		dbAction.Description,
		dbAction.Parent,
		dbAction.Enabled,
		dbAction.UpdatedAt,
		&dbAction.DeletedAt.Time,
	), nil

}

func (r *actionRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Action, error) {
	var dbActions []Action

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	err := r.db.Model(&Action{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).
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
			dbAction.TenantID,
			dbAction.BranchID,
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

func (r *actionRepository) Remove(ctx context.Context, input *baseCmd.BaseInput, actionID string) error {

	resp := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", actionID).Delete(&Action{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *actionRepository) HasChildren(ctx context.Context, input *baseCmd.BaseInput, actionID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Action{}).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND parent = ?", actionID).
		Count(&count).Error

	return count > 0, err
}
