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

type widgetRepository struct {
	db *gorm.DB
}

func NewWidgetRepository(db *gorm.DB) domain.WidgetRepository {
	return &widgetRepository{db: db}
}

func (r *widgetRepository) Create(ctx context.Context, input *domain.Widget) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}

	return db.WithContext(ctx).Create(widgetMapper(input)).Error
}

func (r *widgetRepository) Update(ctx context.Context, input *domain.Widget) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Save(widgetMapper(input)).Error
}

func (r *widgetRepository) FindByID(ctx context.Context, dashboardID, widgetID string) (*domain.Widget, error) {
	var dbWidget WidgetDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND dashboard_id = ? AND id =?", dashboardID, widgetID).
		First(&dbWidget).Error
	if err != nil {
		return nil, err
	}
	return toWidgetDomain(dbWidget), nil

}

func (r *widgetRepository) FindAll(ctx context.Context, dashboardID string) ([]*domain.Widget, error) {
	var dbWidgets []WidgetDB

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
	err = db.Model(&WidgetDB{}).Where(gormhelper.DeleteFilter()+" AND dashboard_id = ?", dashboardID).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()).
		Preload("Link").
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbWidgets).Error
	if err != nil {
		return nil, err
	}

	return toWidgetsDomain(dbWidgets), nil
}

func (r *widgetRepository) Remove(ctx context.Context, dashboardID, widgetID string) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}

	// Start a transaction
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, delete the related link data
		if err := tx.Where("widget_id = ?", widgetID).Delete(&WidgetLinkData{}).Error; err != nil {
			return err
		}

		// Then, delete the widget
		result := tx.Where(
			gormhelper.DeleteFilter()+" AND dashboard_id = ? AND id = ?", dashboardID, widgetID).Delete(&WidgetDB{})

		if result.RowsAffected == 0 {
			return errors.New("no widget deleted")
		}

		return result.Error
	})
}

func toWidgetsDomain(dbWidgets []WidgetDB) (widgets []*domain.Widget) {
	widgets = make([]*domain.Widget, len(dbWidgets))
	for i, dbWidget := range dbWidgets {
		widgets[i] = toWidgetDomain(dbWidget)
	}
	return
}
func toWidgetDomain(dbWidget WidgetDB) (widgets *domain.Widget) {

	wData := domain.NewWidgetData(dbWidget.Label, dbWidget.ShowLabel, dbWidget.ShowEmotion, dbWidget.TrueEmotion)
	lMap := make([]domain.WidgetLinkData, 0)
	for l := range dbWidget.Link {
		lMap = append(lMap, domain.NewWidgetLinkData(
			dbWidget.Link[l].Measure,
			dbWidget.Link[l].Tag,
			dbWidget.Link[l].Legend,
		))
	}
	wData.SetLink(lMap)
	wData.SetOptions(dbWidget.Options)

	return domain.NewWidgetFromRepository(
		dbWidget.ID,
		"dbWidget.TenantID", "dbWidget.BranchName",
		dbWidget.DashboardID,
		dbWidget.TypeWidget,
		dbWidget.X, dbWidget.Y,
		dbWidget.W, dbWidget.H,
		wData,
		dbWidget.UpdatedAt, dbWidget.DeletedAt)
}
