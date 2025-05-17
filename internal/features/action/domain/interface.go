package domain

import (
	"context"
)

type ActionRepository interface {
	Create(ctx context.Context, input *Action) error
	Update(ctx context.Context, input *Action) error
	FindByID(ctx context.Context, assetID string) (*Action, error)
	FindAll(ctx context.Context) ([]*Action, error)
	Remove(ctx context.Context, actionID string) error
}
