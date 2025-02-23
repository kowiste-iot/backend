package services

import (
	"backend/shared/stream/domain"
	"backend/shared/stream/infra/nats"
	"backend/shared/stream/repo"
	"fmt"
	"time"
)

func (c *Container) initializeStreamService(s *Services) (err error) {

	messageRepo := repository.NewMessageRepository(c.base.DB)
	natsFactory := nats.NewNatsClientFactory(messageRepo)

	// Configure NATS stream
	streamConfig := &domain.StreamConfig{
		URL:            "http://localhost:4222",
		MaxReconnects:  5,
		ReconnectWait:  time.Second * 5,
		ConnectTimeout: time.Second * 10,
		WriteTimeout:   time.Second * 5,
		PersistMessage: true,
	}

	// Create NATS client
	s.StreamService, err = natsFactory.CreateClient(streamConfig)
	if err != nil {
		return fmt.Errorf("failed to create NATS client: %w", err)
	}

	return nil
}
