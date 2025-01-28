package domain

import (
	"context"
	baseCmd "ddd/shared/base/command"
)

type DashboardRepository interface {
	Create(ctx context.Context, input *Dashboard) error
	Update(ctx context.Context, input *Dashboard) error
	FindByID(ctx context.Context, input *baseCmd.BaseInput, assetID string) (*Dashboard, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*Dashboard, error)
	Remove(ctx context.Context, input *baseCmd.BaseInput, dashboardID string) error
}
