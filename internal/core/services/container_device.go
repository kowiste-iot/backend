package services

import (
	appDevice "backend/internal/features/device/app"
	repoDevice "backend/internal/features/device/infra/gorm"
	"errors"
)

func (c *Container) initializeDeviceService(s *Services) error {
	if s.AssetDepService == nil {
		return errors.New("asset dependency service must be initialized first")
	}

	deviceRepo := repoDevice.NewRepository(c.base.DB)
	s.DeviceService = appDevice.NewService(c.base, deviceRepo, s.AssetDepService)
	return nil
}
