package command

import (
    "backend/shared/base/command"
    "time"
)

// MessageIDInput represents input for retrieving a specific message
type MessageIDInput struct {
    command.BaseInput
    MessageID string `json:"messageId" validate:"required"`
}

// TimeRangeInput represents input for time-based message queries
type TimeRangeInput struct {
    command.BaseInput
    StartTime time.Time `json:"startTime" validate:"required"`
    EndTime   time.Time `json:"endTime" validate:"required"`
}

// BatchStoreInput represents input for storing multiple messages
type BatchStoreInput struct {
    command.BaseInput
    Messages []MessageData `json:"messages" validate:"required,min=1"`
}

// MessageData represents the data structure for message storage
type MessageData struct {
    ID       string                 `json:"id" validate:"required"`
    TenantID string                 `json:"tenantId" validate:"required"`
    BranchID string                 `json:"branchId" validate:"required"`
    Time     time.Time              `json:"time" validate:"required"`
    Data     map[string]interface{} `json:"data" validate:"required"`
}