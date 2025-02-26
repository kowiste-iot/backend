package domain

import (
	"backend/internal/features/datastore/domain/command"
	ingestDomain "backend/internal/features/ingest/domain"
	baseCmd "backend/shared/base/command"

	"context"
)

// MessageRepository defines the interface for message storage operations
type MessageRepository interface {
	// Store saves a single message
	Store(ctx context.Context, msg *ingestDomain.Message) error

	// StoreBatch saves multiple messages in a batch
	StoreBatch(ctx context.Context, msgs []*ingestDomain.Message) error

	// FindByID retrieves a message by its ID
	FindByID(ctx context.Context, input *baseCmd.BaseInput, id string) (*ingestDomain.Message, error)

	// FindAll retrieves all messages for a tenant/branch
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*ingestDomain.Message, error)

	// FindByTimeRange retrieves messages within a time range
	FindByTimeRange(ctx context.Context, input *command.TimeRangeInput) ([]*ingestDomain.Message, error)
}
