package services

import (
	appDevice "backend/internal/features/device/app"
	emqx "backend/internal/features/device/infra/emqx/rest"
	repoDevice "backend/internal/features/device/infra/gorm"
	"errors"
)

func (c *Container) initializeDeviceService(s *Services) error {
	if s.AssetDepService == nil {
		return errors.New("asset dependency service must be initialized first")
	}

	deviceRepo := repoDevice.NewRepository(c.base.DB)
	deviceBroker := emqx.NewDeviceBroker(&c.config.Ingest)
	s.DeviceService = appDevice.NewService(c.base, &appDevice.ServiceDependencies{
		Repo:     deviceRepo,
		AssetDep: s.AssetDepService,
		Broker:   deviceBroker,
	})
	return nil
}
