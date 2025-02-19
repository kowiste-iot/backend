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

type widgetRepository struct {
	db *gorm.DB
}

func NewWidgetRepository(db *gorm.DB) domain.WidgetRepository {
	db.AutoMigrate(&WidgetLinkData{}, &Widget{})
	return &widgetRepository{db: db}
}

func (r *widgetRepository) Create(ctx context.Context, input *domain.Widget) error {
	dbWidget := Widget{
		ID:         input.ID(),
		TenantID:   input.TenantID(),
		BranchName: input.BranchName(),
		Name:       input.Name(),
	}
	return r.db.WithContext(ctx).Create(&dbWidget).Error
}

func (r *widgetRepository) Update(ctx context.Context, input *domain.Widget) error {
	dbWidget := Widget{
		ID:         input.ID(),
		TenantID:   input.TenantID(),
		BranchName: input.BranchName(),
		Name:       input.Name(),
	}
	return r.db.WithContext(ctx).Updates(&dbWidget).Error
}

func (r *widgetRepository) FindByID(ctx context.Context, input *baseCmd.BaseInput, dashboardID, widgetID string) (*domain.Widget, error) {
	var dbWidget Widget

	err := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", dashboardID).
		First(&dbWidget).Error
	if err != nil {
		return nil, err
	}
	return domain.NewWidgetFromRepository(
		dbWidget.ID,
		dbWidget.TenantID,
		dbWidget.BranchName,
		dbWidget.Name,
		"",
		0,
		0,
		0,
		0,
		0,
		0,
		domain.WidgetData{},
		time.Now(),
		nil,
	), nil

}

func (r *widgetRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Widget, error) {
	var dbWidgets []Widget

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	err := r.db.Model(&Widget{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbWidgets).Error
	if err != nil {
		return nil, err
	}

	widgets := make([]*domain.Widget, len(dbWidgets))
	for i, dbWidget := range dbWidgets {
		widgets[i] = domain.NewWidgetFromRepository(
			dbWidget.ID,
			dbWidget.TenantID,
			dbWidget.BranchName,
			dbWidget.Name,
			"",
			0,
			0,
			0,
			0,
			0,
			0,
			domain.WidgetData{},
			time.Now(),
			nil,
		)
	}
	return widgets, nil
}

func (r *widgetRepository) Remove(ctx context.Context, input *baseCmd.BaseInput, dashboardID, widgetID string) error {

	resp := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND dashboard_id = ? AND id = ?", dashboardID, widgetID).Delete(&Widget{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}
