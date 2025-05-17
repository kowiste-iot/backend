package command

import (
	scopeDomain "backend/internal/features/scope/domain"
	"backend/shared/base/command"

	"fmt"
)

type CreatePermissionInput struct {
	command.BaseInput
	Name             string
	Description      string
	Type             string
	ResourceType     string
	Resources        string
	Scopes           []string
	Policies         []string
	DecisionStrategy string
}
type UpdatePermissionInput struct {
	command.BaseInput
	ID           string
	ResourceID   string
	ResourceName string
	Roles        map[string][]scopeDomain.Scope
}

func ResourceName(roleName string) string {
	return fmt.Sprintf("%s-resource", roleName)
}
