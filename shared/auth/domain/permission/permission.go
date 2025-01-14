package permission

import "context"

const (
	TypeScope string = "scope"
	TypeRole  string = "role"
)
const (
	DecisionUnanimous   string = "UNANIMOUS"
	DecisionAffirmative string = "AFFIRMATIVE"
)
const (
	LogicPositive string = "POSITIVE"
)

type Permission struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Type             string   `json:"type"` // resource-based, scope-based
	Resources        []string `json:"resources,omitempty"`
	Scopes           []string `json:"scopes,omitempty"`
	Policies         []string `json:"policies"`
	DecisionStrategy string   `json:"decisionStrategy"` // UNANIMOUS, AFFIRMATIVE, CONSENSUS
}

type PermissionProvider interface {
	CreatePermission(ctx context.Context, tenantID, clientID string, permission Permission) (*Permission, error)
	UpdatePermission(ctx context.Context, tenantID, clientID string, permission Permission) error
	DeletePermission(ctx context.Context, tenantID, clientID, permissionID string) error
	GetPermission(ctx context.Context, tenantID, clientID, permissionID string) (*Permission, error)
	ListPermissions(ctx context.Context, tenantID, clientID string) ([]Permission, error)
}
