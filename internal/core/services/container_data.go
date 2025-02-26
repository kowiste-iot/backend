package services

import (
	"backend/internal/features/datastore/infra/gorm"
	"backend/internal/features/datastore/app"
	"time"
)

func (c *Container) initializeMessageStoreService(s *Services) error {
	repo := gorm.NewRepository(c.base.DB)

	config := app.ServiceConfig{
		BatchSize:     100,
		FlushInterval: 5 * time.Second,
	}

	s.DataStoreService = app.NewService(
		c.base,
		repo,
		s.StreamService,
		config,
	)

	return nil
}
