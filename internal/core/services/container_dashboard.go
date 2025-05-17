package services

import (
	appDashboard "backend/internal/features/dashboard/app"
	repoDashboard "backend/internal/features/dashboard/infra/gorm"
	"errors"
)

func (c *Container) initializeDashboardService(s *Services) error {
	if s.AssetDepService == nil {
		return errors.New("asset dependency service must be initialized first")
	}

	dashboardRepo := repoDashboard.NewRepository(c.base.DB)
	s.DashboardService = appDashboard.NewService(c.base, dashboardRepo, s.AssetDepService)

	widgetRepo := repoDashboard.NewWidgetRepository(c.base.DB)
	s.WidgetService = appDashboard.NewWidgetService(c.base, widgetRepo)

	return nil
}
