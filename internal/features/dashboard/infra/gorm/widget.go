package repository

import (
	"backend/internal/features/dashboard/domain"
	baseCmd "backend/shared/base/command"
	gormhelper "backend/shared/gorm"
	"backend/shared/pagination"
	"context"
	"errors"

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
	lMap := make([]WidgetLinkData, 0)
	for _, link := range input.Link() {
		lMap = append(lMap, WidgetLinkData{
			WidgetID: input.ID(),
			Measure:  link.MeasureID(),
			Tag:      link.Tag(),
			Legend:   link.Legend(),
		})
	}
	dbWidget := Widget{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchName:  input.BranchName(),
		DashboardID: input.DashboardID(),
		TypeWidget:  input.TypeWidget(),
		X:           input.X(),
		Y:           input.Y(),
		W:           input.W(),
		H:           input.H(),
		Label:       input.Label(),
		ShowLabel:   input.ShowLabel(),
		ShowEmotion: input.ShowEmotion(),
		TrueEmotion: input.TrueEmotion(),
		Link:        lMap,
	}
	return r.db.WithContext(ctx).Create(&dbWidget).Error
}

func (r *widgetRepository) Update(ctx context.Context, input *domain.Widget) error {
	dbWidget := Widget{
		ID:         input.ID(),
		TenantID:   input.TenantID(),
		BranchName: input.BranchName(),
	}
	return r.db.WithContext(ctx).Updates(&dbWidget).Error
}

func (r *widgetRepository) FindByID(ctx context.Context, input *baseCmd.BaseInput, dashboardID, widgetID string) (*domain.Widget, error) {
	var dbWidget Widget

	err := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND dashboard_id = ?", dashboardID).
		First(&dbWidget).Error
	if err != nil {
		return nil, err
	}
	return toWidgetDomain(dbWidget), nil

}

func (r *widgetRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput, dashboardID string) ([]*domain.Widget, error) {
	var dbWidgets []Widget

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	err := r.db.Model(&Widget{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND dashboard_id = ?", dashboardID).Count(&total).Error
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

	return toWidgetsDomain(dbWidgets), nil
}

func (r *widgetRepository) Remove(ctx context.Context, input *baseCmd.BaseInput, dashboardID, widgetID string) error {
	resp := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND dashboard_id = ? AND id = ?", dashboardID, widgetID).Delete(&Widget{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func toWidgetsDomain(dbWidgets []Widget) (widgets []*domain.Widget) {
	widgets = make([]*domain.Widget, len(dbWidgets))
	for i, dbWidget := range dbWidgets {
		widgets[i] = toWidgetDomain(dbWidget)
	}
	return
}
func toWidgetDomain(dbWidget Widget) (widgets *domain.Widget) {

	wData := domain.NewWidgetData(dbWidget.Label, dbWidget.ShowLabel, dbWidget.ShowEmotion, dbWidget.TrueEmotion)
	lMap := make([]domain.WidgetLinkData, 0)
	for l := range dbWidget.Link {
		lMap = append(lMap, domain.NewWidgetLinkData(
			dbWidget.Link[l].Measure,
			dbWidget.Link[l].Tag,
			dbWidget.Link[l].Measure,
		))
	}
	wData.SetLink(lMap)
	wData.SetOptions(dbWidget.Options)

	return domain.NewWidgetFromRepository(
		dbWidget.ID,
		dbWidget.TenantID, dbWidget.BranchName,
		dbWidget.DashboardID,
		dbWidget.TypeWidget,
		dbWidget.X, dbWidget.Y,
		dbWidget.W, dbWidget.H,
		wData,
		dbWidget.UpdatedAt, dbWidget.DeletedAt)
}
