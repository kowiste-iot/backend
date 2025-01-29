package repository

import (
	"context"
	"ddd/internal/features/alert/domain"
	baseCmd "ddd/shared/base/command"
	gormhelper "ddd/shared/gorm"
	"ddd/shared/pagination"
	"errors"
	"time"

	"gorm.io/gorm"
)

type alertRepository struct {
	db *gorm.DB
}

type Alert struct {
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

func NewRepository(db *gorm.DB) domain.AlertRepository {
	db.AutoMigrate(&Alert{})
	return &alertRepository{db: db}
}

func (r *alertRepository) Create(ctx context.Context, input *domain.Alert) error {
	dbAlert := Alert{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Create(&dbAlert).Error
}

func (r *alertRepository) Update(ctx context.Context, input *domain.Alert) error {
	dbAlert := Alert{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Updates(&dbAlert).Error
}

func (r *alertRepository) FindByID(ctx context.Context, input *baseCmd.BaseInput, alertID string) (*domain.Alert, error) {
	var dbAlert Alert

	err := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", alertID).
		First(&dbAlert).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbAlert.ID,
		dbAlert.TenantID,
		dbAlert.BranchID,
		dbAlert.Name,
		dbAlert.Description,
		dbAlert.Parent,
		dbAlert.Enabled,
		dbAlert.UpdatedAt,
		&dbAlert.DeletedAt.Time,
	), nil

}

func (r *alertRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Alert, error) {
	var dbAlerts []Alert

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	err := r.db.Model(&Alert{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbAlerts).Error
	if err != nil {
		return nil, err
	}

	alerts := make([]*domain.Alert, len(dbAlerts))
	for i, dbAlert := range dbAlerts {
		alerts[i] = domain.NewFromRepository(
			dbAlert.ID,
			dbAlert.TenantID,
			dbAlert.BranchID,
			dbAlert.Name,
			dbAlert.Description,
			dbAlert.Parent,
			dbAlert.Enabled,
			dbAlert.UpdatedAt,
			&dbAlert.DeletedAt.Time,
		)
	}
	return alerts, nil
}

func (r *alertRepository) Remove(ctx context.Context, input *baseCmd.BaseInput, alertID string) error {

	resp := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", alertID).Delete(&Alert{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *alertRepository) HasChildren(ctx context.Context, input *baseCmd.BaseInput, alertID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Alert{}).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND parent = ?", alertID).
		Count(&count).Error

	return count > 0, err
}
