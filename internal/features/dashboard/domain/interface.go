package domain

import (
	baseCmd "backend/shared/base/command"
	"context"
)

type DashboardRepository interface {
	Create(ctx context.Context, input *Dashboard) error
	Update(ctx context.Context, input *Dashboard) error
	FindByID(ctx context.Context, input *baseCmd.BaseInput, dashboardID string) (*Dashboard, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*Dashboard, error)
	Remove(ctx context.Context, input *baseCmd.BaseInput, dashboardID string) error
}
type WidgetRepository interface {
	Create(ctx context.Context, input *Widget) error
	Update(ctx context.Context, input *Widget) error
	FindByID(ctx context.Context, input *baseCmd.BaseInput, dashboardID, widgetID string) (*Widget, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput, dashboardID string) ([]*Widget, error)
	Remove(ctx context.Context, input *baseCmd.BaseInput, dashboardID, widgetID string) error
}
