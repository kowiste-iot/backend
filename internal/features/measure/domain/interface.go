package domain

import (
	"context"
	baseCmd "ddd/shared/base/command"
)

type MeasureRepository interface {
	Create(ctx context.Context, input *Measure) error
	Update(ctx context.Context, input *Measure) error
	FindByID(ctx context.Context, input *baseCmd.BaseInput, assetID string) (*Measure, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*Measure, error)
	Remove(ctx context.Context, input *baseCmd.BaseInput, measureID string) error
}
