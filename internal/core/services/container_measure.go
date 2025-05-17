package services

import (
	appMeasure "backend/internal/features/measure/app"
	repoMeasure "backend/internal/features/measure/infra/gorm"
	"errors"
)

func (c *Container) initializeMeasureService(s *Services) error {
	if s.AssetDepService == nil {
		return errors.New("asset dependency service must be initialized first")
	}

	measureRepo := repoMeasure.NewRepository(c.base.DB)
	s.MeasureService = appMeasure.NewService(c.base, measureRepo, s.AssetDepService)
	return nil
}
