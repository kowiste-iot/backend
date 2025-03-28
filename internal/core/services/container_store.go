package services

import (
	"backend/internal/features/datastore/app"
	indestDomain "backend/internal/features/ingest/domain"

	repository "backend/internal/features/datastore/infra/gorm"
	"context"
	"time"
)

func (c *Container) initializeStoreService(s *Services) (err error) {
	config := &app.ServiceConfig{
		BatchSize:     100,
		FlushInterval: 5 * time.Second,
		TopicNats:     indestDomain.TopicIngest,
	}

	stream, err := c.initializeStreamService(s)
	if err != nil {
		return
	}
	storeRepo := repository.NewRepository(c.base.DB)
	s.DataStoreService = app.NewService(c.base, storeRepo, stream, config)
	err = s.DataStoreService.Start(context.Background())

	return
}
