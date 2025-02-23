package domain

import "time"

type Message struct {
	ID       string                 `json:"id"`
	TenantID string                 `json:"tenant"`
	BranchID string                 `json:"branch"`
	Time     time.Time              `json:"time"`
	Data     map[string]interface{} `json:"data"`
}
