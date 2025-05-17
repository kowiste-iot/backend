package services

import (
	appAction "backend/internal/features/action/app"
	repoAction "backend/internal/features/action/infra/gorm"
	"errors"
)

func (c *Container) initializeActionService(s *Services) error {
	if s.AssetDepService == nil {
		return errors.New("asset dependency service must be initialized first")
	}

	actionRepo := repoAction.NewRepository(c.base.DB)
	s.ActionService = appAction.NewService(c.base, actionRepo, s.AssetDepService)
	return nil
}
