package domain

import (
	baseCmd "backend/shared/base/command"
	"context"
)

type AssetRepository interface {
	Create(ctx context.Context, input *Asset) error
	Update(ctx context.Context, input *Asset) error
	FindByID(ctx context.Context, input *baseCmd.BaseInput, assetID string) (*Asset, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*Asset, error)
	Remove(ctx context.Context, input *baseCmd.BaseInput, assetID string) error
	HasChildren(ctx context.Context, input *baseCmd.BaseInput, assetID string) (bool, error)
}
