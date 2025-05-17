package domain

import (
	"context"
)

type DeviceRepository interface {
	Create(ctx context.Context, input *Device) (password string, err error)
	Update(ctx context.Context, input *Device) error
	FindByID(ctx context.Context, assetID string) (*Device, error)
	FindAll(ctx context.Context) ([]*Device, error)
	Remove(ctx context.Context, deviceID string) error
}

type BrokerProvider interface{
	KickOut(ctx context.Context, clientID string)(err error)
}