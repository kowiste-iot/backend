// domain/store.go
package domain

import (
	"time"
)

type MessageStore struct {
	id        string
	tenantID  string
	branchID  string
	time      time.Time
	data      map[string]interface{}
	createdAt time.Time
}

// ID returns the message ID
func (m *MessageStore) ID() string {
	return m.id
}

// TenantID returns the tenant ID
func (m *MessageStore) TenantID() string {
	return m.tenantID
}

// BranchID returns the branch ID
func (m *MessageStore) BranchID() string {
	return m.branchID
}

// Time returns the message time
func (m *MessageStore) Time() time.Time {
	return m.time
}

// Data returns the message data
func (m *MessageStore) Data() map[string]interface{} {
	return m.data
}

// CreatedAt returns when the message was stored
func (m *MessageStore) CreatedAt() time.Time {
	return m.createdAt
}
