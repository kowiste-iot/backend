package app

import (
	"backend/internal/features/datastore/domain"
	ingestDomain "backend/internal/features/ingest/domain"
	"backend/shared/base"
	streamDomain "backend/shared/stream/domain"
	"backend/shared/validator"
	"context"
	"fmt"
	"sync"
	"time"
)

type DataStoreService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	StoreBatch(ctx context.Context, messages []*ingestDomain.Message) error
}

type messageStoreService struct {
	repo         domain.MessageRepository
	streamClient streamDomain.StreamClient
	*base.BaseService

	buffer   []*ingestDomain.Message
	bufferMu sync.Mutex
	stopChan chan struct{}
	*ServiceConfig
}

type ServiceConfig struct {
	BatchSize     int
	FlushInterval time.Duration
	TopicNats     string
}

func NewService(base *base.BaseService, repo domain.MessageRepository, streamClient streamDomain.StreamClient, config *ServiceConfig) DataStoreService {
	if config.BatchSize == 0 {
		config.BatchSize = 100
	}
	if config.FlushInterval == 0 {
		config.FlushInterval = 5 * time.Second
	}

	return &messageStoreService{
		repo:          repo,
		streamClient:  streamClient,
		BaseService:   base,
		buffer:        make([]*ingestDomain.Message, 0, config.BatchSize),
		ServiceConfig: config,
		stopChan:      make(chan struct{}),
	}
}

func (s *messageStoreService) handleMessage(ctx context.Context, msg *streamDomain.WireMessage) (err error) {
	if err := validator.Validate(msg); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	s.bufferMu.Lock()
	ingData := new(ingestDomain.Message)
	err = msg.DataToModel(ingData)
	if err != nil {
		return
	}
	s.buffer = append(s.buffer, ingData)
	shouldFlush := len(s.buffer) >= s.BatchSize
	s.bufferMu.Unlock()

	if shouldFlush {
		return s.flush(ctx)
	}
	return nil
}

func (s *messageStoreService) flush(ctx context.Context) error {
	s.bufferMu.Lock()
	if len(s.buffer) == 0 {
		s.bufferMu.Unlock()
		return nil
	}

	messages := make([]*ingestDomain.Message, len(s.buffer))
	copy(messages, s.buffer)
	s.buffer = s.buffer[:0]
	s.bufferMu.Unlock()

	return s.StoreBatch(ctx, messages)
}

func (s *messageStoreService) StoreBatch(ctx context.Context, messages []*ingestDomain.Message) error {
	if err := s.repo.StoreBatch(ctx, messages); err != nil {
		return fmt.Errorf("failed to store message batch: %w", err)
	}
	return nil
}

func (s *messageStoreService) startFlushTimer(ctx context.Context) {
	ticker := time.NewTicker(s.FlushInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.flush(context.Background()); err != nil {
					// Log error
					s.Logger.Error(ctx, err, "failed to flush messages", "error")
				}
			case <-s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *messageStoreService) Start(ctx context.Context) error {
	s.startFlushTimer(ctx)
	if err := s.streamClient.Subscribe(ctx, s.TopicNats, s.handleMessage); err != nil {
		return fmt.Errorf("failed to subscribe to messages: %w", err)
	}
	return nil
}

func (s *messageStoreService) Stop(ctx context.Context) error {
	close(s.stopChan)
	if err := s.flush(ctx); err != nil {
		s.Logger.Error(ctx, err, "failed to flush messages during shutdown", "error")
	}
	if err := s.streamClient.Unsubscribe(ctx, s.TopicNats); err != nil {
		return fmt.Errorf("failed to unsubscribe from messages: %w", err)
	}
	return nil
}
