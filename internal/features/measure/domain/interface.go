package domain

import (
	"context"
)

type MeasureRepository interface {
	Create(ctx context.Context, input *Measure) error
	Update(ctx context.Context, input *Measure) error
	FindByID(ctx context.Context, assetID string) (*Measure, error)
	FindAll(ctx context.Context) ([]*Measure, error)
	Remove(ctx context.Context, measureID string) error
}
