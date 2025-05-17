package domain

import (
	"context"
)

type DashboardRepository interface {
	Create(ctx context.Context, input *Dashboard) error
	Update(ctx context.Context, input *Dashboard) error
	FindByID(ctx context.Context, dashboardID string) (*Dashboard, error)
	FindAll(ctx context.Context) ([]*Dashboard, error)
	Remove(ctx context.Context, dashboardID string) error
}
type WidgetRepository interface {
	Create(ctx context.Context, input *Widget) error
	Update(ctx context.Context, input *Widget) error
	FindByID(ctx context.Context, dashboardID, widgetID string) (*Widget, error)
	FindAll(ctx context.Context, dashboardID string) ([]*Widget, error)
	Remove(ctx context.Context, dashboardID, widgetID string) error
}
