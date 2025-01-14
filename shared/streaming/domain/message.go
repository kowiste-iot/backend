package domain

import "encoding/json"

type Message struct {
	TenantID  string `json:"tenant_id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
	Priority  string `json:"priority"`
	Data      any    `json:"data"`
	CreatedAt int64  `json:"created_at"`
}

func (m Message) ToByte() (data []byte, err error) {
	return json.Marshal(m)
}
