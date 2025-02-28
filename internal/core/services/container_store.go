package services

import (
	"backend/internal/features/datastore/app"
	repository "backend/internal/features/datastore/infra/gorm"
	"context"
	"time"
)

func (c *Container) initializeStoreService(s *Services) (err error) {
	config := &app.ServiceConfig{
		BatchSize:     100,
		FlushInterval: 5 * time.Second,
		TopicNats:     "data.ingest",
	}
	storeRepo := repository.NewRepository(c.base.DB)
	s.DataStoreService = app.NewService(c.base, storeRepo, s.StreamService, config)
	err = s.DataStoreService.Start(context.Background())

	return
}
