package domain

import "errors"

var (
    ErrMessageNotFound     = errors.New("message not found")
    ErrInvalidTimeRange    = errors.New("invalid time range")
    ErrInvalidMessageData  = errors.New("invalid message data")
    ErrBatchStoreFailed    = errors.New("failed to store message batch")
    ErrInvalidBatchSize    = errors.New("invalid batch size")
)