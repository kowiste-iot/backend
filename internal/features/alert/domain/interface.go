package domain

import (
	"context"
	baseCmd "ddd/shared/base/command"
)

type AlertRepository interface {
	Create(ctx context.Context, input *Alert) error
	Update(ctx context.Context, input *Alert) error
	FindByID(ctx context.Context, input *baseCmd.BaseInput, assetID string) (*Alert, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*Alert, error)
	Remove(ctx context.Context, input *baseCmd.BaseInput, alertID string) error
}
