package domain

import (
	"context"
	"ddd/internal/features/asset/domain/command"
	baseCmd "ddd/shared/base/command"
)

type AssetRepository interface {
	Create(ctx context.Context, input *command.CreateAssetInput) error
	Update(ctx context.Context, input *command.UpdateAssetInput) error
	FindByID(ctx context.Context, input *command.AssetIDInput) (*Asset, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*Asset, error)
	Remove(ctx context.Context, input *command.AssetIDInput) error
	HasChildren(ctx context.Context, input *command.AssetIDInput) (bool, error)
}
