package domain

import (
	baseCmd "backend/shared/base/command"
	"context"
)

type ActionRepository interface {
	Create(ctx context.Context, input *Action) error
	Update(ctx context.Context, input *Action) error
	FindByID(ctx context.Context, input *baseCmd.BaseInput, assetID string) (*Action, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*Action, error)
	Remove(ctx context.Context, input *baseCmd.BaseInput, actionID string) error
}
