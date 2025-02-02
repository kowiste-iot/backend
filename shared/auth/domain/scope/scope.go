package scope

import (
	baseCmd "backend/shared/base/command"
	"context"
)

const (
	View   string = "view"
	Create string = "create"
	Update string = "update"
	Delete string = "delete"
)

type Scope struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

type ScopeProvider interface {
	CreateScope(ctx context.Context, tenantID, clientID string, scope Scope) (*Scope, error)
	ListScopes(ctx context.Context, input *baseCmd.BaseInput) ([]Scope, error)
}

func AllScopes() []string {
	return []string{View, Create, Update, Delete}
}
