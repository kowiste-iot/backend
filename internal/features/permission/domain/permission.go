package domain

import (
	roleDomain "backend/internal/features/role/domain"
	scopeDomain "backend/internal/features/scope/domain"
	baseCmd "backend/shared/base/command"

	"context"
	"fmt"

	"github.com/google/uuid"
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

type PermissionProvider interface {
	CreatePermission(ctx context.Context,scopes []scopeDomain.Scope, input *baseCmd.BaseInput, permission *Permission) (*Permission, error)
	ListPermissions(ctx context.Context, input *baseCmd.BaseInput) ([]Permission, error)
}

type Permission struct {
	ID               string            `json:"id,omitempty"`
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	Type             string            `json:"type"` // resource-based, scope-based
	ResourceType     string            `json:"resourceType"`
	Resources        []string          `json:"resources,omitempty"`
	Scopes           []string          `json:"scopes,omitempty"`
	Policies         []string          `json:"policies"`
	Roles            []roleDomain.Role `json:"roles,omitempty"`
	DecisionStrategy string            `json:"decisionStrategy"` // UNANIMOUS, AFFIRMATIVE, CONSENSUS
	Logic            string            `json:"logic"`
}

func New(name, description, typePermission, decisionStrategy string, resources, scopes, policies []string) (resource *Permission, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}

	resource = &Permission{
		ID:               id.String(),
		Name:             name,
		Description:      description,
		Type:             typePermission,
		Resources:        resources,
		Scopes:           scopes,
		Policies:         policies,
		DecisionStrategy: decisionStrategy,
	}
	return
}

func NameNonAdmin(roleName, resourceName string) string {
	return fmt.Sprintf("%s-%s-permission", roleName, resourceName)
}
func NameAdmin() string {
	return fmt.Sprintf("%s-permission", roleDomain.RoleAdmin)
}
