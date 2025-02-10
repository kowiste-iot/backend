package domain

import (scopeDomain "backend/internal/features/scope/domain")

type ResourcePermission struct {
	ID          string                   `json:"id,omitempty"`
	Name        string                   `json:"name"`
	DisplayName string                   `json:"displayName,omitempty"`
	Roles       map[string][]scopeDomain.Scope `json:"roles,omitempty"`
}
