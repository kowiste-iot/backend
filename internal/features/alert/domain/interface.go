package domain

import (
	"context"
)

type AlertRepository interface {
	Create(ctx context.Context, input *Alert) error
	Update(ctx context.Context, input *Alert) error
	FindByID(ctx context.Context, assetID string) (*Alert, error)
	FindAll(ctx context.Context) ([]*Alert, error)
	Remove(ctx context.Context, alertID string) error
}
