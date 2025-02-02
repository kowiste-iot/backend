package resource

import "backend/shared/auth/domain/scope"

type ResourcePermission struct {
	ID          string                   `json:"id,omitempty"`
	Name        string                   `json:"name"`
	DisplayName string                   `json:"displayName,omitempty"`
	Roles       map[string][]scope.Scope `json:"roles,omitempty"`
}
