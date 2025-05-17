package repository

import (
	"backend/internal/features/alert/domain"
	gormhelper "backend/shared/gorm"
	"backend/shared/pagination"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type alertRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.AlertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) Create(ctx context.Context, input *domain.Alert) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Create(alertMapper(input)).Error
}

func (r *alertRepository) Update(ctx context.Context, input *domain.Alert) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Updates(alertMapper(input)).Error
}

func (r *alertRepository) FindByID(ctx context.Context, alertID string) (*domain.Alert, error) {
	var dbAlert AlertDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", alertID).
		First(&dbAlert).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbAlert.ID,
		"dbAlert.TenantID",
		"dbAlert.BranchID",
		dbAlert.Name,
		dbAlert.Description,
		dbAlert.Parent,
		dbAlert.Enabled,
		dbAlert.UpdatedAt,
		&dbAlert.DeletedAt.Time,
	), nil

}

func (r *alertRepository) FindAll(ctx context.Context) ([]*domain.Alert, error) {
	var dbAlerts []AlertDB

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
	err = db.Model(&AlertDB{}).Where(gormhelper.DeleteFilter()).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()).
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
			"dbAlert.TenantID",
			"dbAlert.BranchID",
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

func (r *alertRepository) Remove(ctx context.Context, alertID string) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	resp := db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", alertID).Delete(&AlertDB{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}
