package repository

import (
	"backend/internal/features/dashboard/domain"
	baseCmd "backend/shared/base/command"
	gormhelper "backend/shared/gorm"
	"backend/shared/pagination"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

type Dashboard struct {
	ID          string `gorm:"primaryKey"`
	TenantID    string `gorm:"index"`
	BranchID    string `gorm:"index"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func NewRepository(db *gorm.DB) domain.DashboardRepository {
	db.AutoMigrate(&Dashboard{})
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) Create(ctx context.Context, input *domain.Dashboard) error {
	dbDashboard := Dashboard{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Create(&dbDashboard).Error
}

func (r *dashboardRepository) Update(ctx context.Context, input *domain.Dashboard) error {
	dbDashboard := Dashboard{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Updates(&dbDashboard).Error
}

func (r *dashboardRepository) FindByID(ctx context.Context, input *baseCmd.BaseInput, dashboardID string) (*domain.Dashboard, error) {
	var dbDashboard Dashboard

	err := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", dashboardID).
		First(&dbDashboard).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbDashboard.ID,
		dbDashboard.TenantID,
		dbDashboard.BranchID,
		dbDashboard.Name,
		dbDashboard.Description,
		dbDashboard.Parent,
		dbDashboard.UpdatedAt,
		&dbDashboard.DeletedAt.Time,
	), nil

}

func (r *dashboardRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Dashboard, error) {
	var dbDashboards []Dashboard

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	err := r.db.Model(&Dashboard{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).
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
			dbDashboard.TenantID,
			dbDashboard.BranchID,
			dbDashboard.Name,
			dbDashboard.Description,
			dbDashboard.Parent,
			dbDashboard.UpdatedAt,
			&dbDashboard.DeletedAt.Time,
		)
	}
	return dashboards, nil
}

func (r *dashboardRepository) Remove(ctx context.Context, input *baseCmd.BaseInput, dashboardID string) error {

	resp := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", dashboardID).Delete(&Dashboard{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *dashboardRepository) HasChildren(ctx context.Context, input *baseCmd.BaseInput, dashboardID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Dashboard{}).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND parent = ?", dashboardID).
		Count(&count).Error

	return count > 0, err
}
