package app

import (
	"backend/internal/features/ingest/domain"
	stream "backend/shared/stream/domain"
	"context"
	"fmt"
)

type IngestService interface {
	Start() error
	Stop() error
	AddConsumer(consumer domain.Consumer)
}
type Service struct {
	consumers    []domain.Consumer
	streamClient stream.StreamClient
	config       *ServiceConfig
}

type ServiceConfig struct {
	Topic           string // NATS topic to publish to
	TenantID        string // Default tenant ID if not in message
	PersistMessages bool
}

func NewService(streamClient stream.StreamClient, config *ServiceConfig) *Service {
	return &Service{
		consumers:    make([]domain.Consumer, 0),
		streamClient: streamClient,
		config:       config,
	}
}

func (s *Service) AddConsumer(consumer domain.Consumer) {
	s.consumers = append(s.consumers, consumer)
}

func (s *Service) Start() error {
	if !s.streamClient.IsConnected() {
		if err := s.streamClient.Connect(); err != nil {
			return fmt.Errorf("failed to connect to stream client: %w", err)
		}
	}

	for _, consumer := range s.consumers {
		if err := consumer.Subscribe(s.handleMessage); err != nil {
			return fmt.Errorf("failed to subscribe consumer: %w", err)
		}

		if err := consumer.Start(); err != nil {
			return fmt.Errorf("failed to start consumer: %w", err)
		}
	}
	return nil
}

func (s *Service) Stop() error {
	for _, consumer := range s.consumers {
		if err := consumer.Stop(); err != nil {
			return fmt.Errorf("failed to stop consumer: %w", err)
		}
	}
	return s.streamClient.Close()
}

func (s *Service) handleMessage(msg *domain.Message) (err error) {
	// Convert input.Message to stream.Message
	streamMsg, err := stream.NewFromIngest(msg)
	if err != nil {
		return
	}

	return s.streamClient.Publish(context.Background(), streamMsg)
}
