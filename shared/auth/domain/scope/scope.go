package scope

import "context"

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
	IconURI     string `json:"iconUri,omitempty"`
}

type ScopeProvider interface {
	CreateScope(ctx context.Context, tenantID, clientID string, scope Scope) (*Scope, error)
	UpdateScope(ctx context.Context, tenantID, clientID string, scope Scope) error
	DeleteScope(ctx context.Context, tenantID, clientID, scopeID string) error
	GetScope(ctx context.Context, tenantID, clientID, scopeID string) (*Scope, error)
	ListScopes(ctx context.Context, tenantID, clientID string) ([]Scope, error)
}
