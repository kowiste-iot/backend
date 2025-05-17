package repository

import (
	"backend/internal/features/dashboard/domain"
	gormhelper "backend/shared/gorm"
	"backend/shared/pagination"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) Create(ctx context.Context, input *domain.Dashboard) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}

	return db.WithContext(ctx).Create(dashboardMapper(input)).Error
}

func (r *dashboardRepository) Update(ctx context.Context, input *domain.Dashboard) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}

	return db.WithContext(ctx).Updates(dashboardMapper(input)).Error
}

func (r *dashboardRepository) FindByID(ctx context.Context, dashboardID string) (*domain.Dashboard, error) {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	var dbDashboard DashboardDB

	err = db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", dashboardID).
		First(&dbDashboard).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbDashboard.ID,
		"dbDashboard.TenantID",
		"dbDashboard.BranchID",
		dbDashboard.Name,
		dbDashboard.Description,
		dbDashboard.Parent,
		dbDashboard.UpdatedAt,
		&dbDashboard.DeletedAt.Time,
	), nil

}

func (r *dashboardRepository) FindAll(ctx context.Context) ([]*domain.Dashboard, error) {
	var dbDashboards []DashboardDB

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
	err = db.Model(&DashboardDB{}).Where(gormhelper.DeleteFilter()).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbDashboards).Error
	if err != nil {
		return nil, err
	}

	dashboards := make([]*domain.Dashboard, len(dbDashboards))
	for i, dbDashboard := range dbDashboards {
		dashboards[i] = domain.NewFromRepository(
			dbDashboard.ID,
			"dbDashboard.TenantID",
			"dbDashboard.BranchID",
			dbDashboard.Name,
			dbDashboard.Description,
			dbDashboard.Parent,
			dbDashboard.UpdatedAt,
			&dbDashboard.DeletedAt.Time,
		)
	}
	return dashboards, nil
}

func (r *dashboardRepository) Remove(ctx context.Context, dashboardID string) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	resp := db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", dashboardID).Delete(&DashboardDB{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}
