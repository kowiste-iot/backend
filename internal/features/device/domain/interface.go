package domain

import (
	"context"
	baseCmd "ddd/shared/base/command"
)

type DeviceRepository interface {
	Create(ctx context.Context, input *Device) error
	Update(ctx context.Context, input *Device) error
	FindByID(ctx context.Context, input *baseCmd.BaseInput, assetID string) (*Device, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*Device, error)
	Remove(ctx context.Context, input *baseCmd.BaseInput, deviceID string) error
}
