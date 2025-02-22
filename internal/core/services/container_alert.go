package services

import (
	appAlert "backend/internal/features/alert/app"
	repoAlert "backend/internal/features/alert/infra/gorm"
	"errors"
)

func (c *Container) initializeAlertService(s *Services) error {
	if s.AssetDepService == nil {
		return errors.New("asset dependency service must be initialized first")
	}

	alertRepo := repoAlert.NewRepository(c.base.DB)
	s.AlertService = appAlert.NewService(c.base, alertRepo, s.AssetDepService)
	return nil
}