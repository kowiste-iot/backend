package domain

import (
	"context"
)

type AssetRepository interface {
	Create(ctx context.Context, input *Asset) error
	Update(ctx context.Context, input *Asset) error
	FindByID(ctx context.Context, assetID string) (*Asset, error)
	FindAll(ctx context.Context) ([]*Asset, error)
	Remove(ctx context.Context, assetID string) error
	HasChildren(ctx context.Context, assetID string) (bool, error)
}
