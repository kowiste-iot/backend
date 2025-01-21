package permission

import (
	"context"
	"ddd/shared/auth/domain/command"
	baseCmd "ddd/shared/base/command"
)

const (
	TypeScope    string = "scope"
	TypeResource string = "resource"
)
const (
	DecisionUnanimous   string = "UNANIMOUS"
	DecisionAffirmative string = "AFFIRMATIVE"
)
const (
	LogicPositive string = "POSITIVE"
)

const (
	defaultPermission string = "Default Permission"
	adminPermission   string = "admin-permission"
)

type Permission struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Type             string   `json:"type"` // resource-based, scope-based
	ResourceType     string   `json:"resourceType"`
	Resources        []string `json:"resources,omitempty"`
	Scopes           []string `json:"scopes,omitempty"`
	Policies         []string `json:"policies"`
	Roles            []string `json:"roles,omitempty"`
	DecisionStrategy string   `json:"decisionStrategy"` // UNANIMOUS, AFFIRMATIVE, CONSENSUS
	Logic            string   `json:"logic"`
}

type PermissionProvider interface {
	CreatePermission(ctx context.Context, tenantID, clientID string, permission Permission) (*Permission, error)
	UpdatePermission(ctx context.Context, tenantID, clientID string, permission Permission) error
	DeletePermission(ctx context.Context, tenantID, clientID, permissionID string) error
	GetPermission(ctx context.Context, tenantID, clientID, permissionID string) (*Permission, error)
	ListPermissions(ctx context.Context, input *baseCmd.BaseInput) ([]Permission, error)
}

type Permissions []Permission

func (rs Permissions) Filter(filterAdmin bool) (permission []Permission) {
	for i := range rs {
		if rs[i].Name == defaultPermission ||
			(filterAdmin && rs[i].Name == adminPermission) {
			continue
		}
		for j := range rs[i].Policies {
			rs[i].Roles = append(rs[i].Roles, command.PolicyToRole(rs[i].Policies[j]))
		}
		permission = append(permission, rs[i])
	}
	return
}
