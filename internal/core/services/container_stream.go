package services

import (
	"backend/shared/stream/domain"
	"backend/shared/stream/infra/nats"
	repository "backend/shared/stream/repo"
	"fmt"
	"time"
)

func (c *Container) initializeStreamService(s *Services) (client domain.StreamClient, err error) {

	messageRepo := repository.NewMessageRepository(c.base.DB)
	natsFactory := nats.NewNatsClientFactory(messageRepo)

	// Configure NATS stream
	streamConfig := &domain.StreamConfig{
		URL:            c.config.Stream.URL,
		MaxReconnects:  5,
		ReconnectWait:  time.Second * 5,
		ConnectTimeout: time.Second * 10,
		WriteTimeout:   time.Second * 5,
		PersistMessage: true,
	}

	// Create NATS client
	client, err = natsFactory.CreateClient(streamConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stream client: %w", err)
	}

	return
}
