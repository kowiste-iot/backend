package command

import (
	"backend/shared/auth/domain/scope"
	"backend/shared/base/command"
	"fmt"
)

type ResourceIDInput struct {
	command.BaseInput
	ResourceID string `validate:"required,uuid"`
}
type ResourceAssignRoleInput struct {
	command.BaseInput
	ResourceID   string `validate:"required,uuid"`
	ResourceName string `validate:"uuid"`
	RoleID       string `validate:"required,uuid"`
	RoleName     string `validate:"uuid"`
	Scopes       []scope.Scope
}

type UpdateResourceInput struct {
	command.BaseInput
	ID          string
	Name        string
	DisplayName string
	Roles       map[string][]scope.Scope
}
type CreateResourceInput struct {
	command.BaseInput
	ID          string
	Name        string
	Type        string
	DisplayName string
	Scopes      []string
}

func ResourceName(roleName string) string {
	return fmt.Sprintf("%s-resource", roleName)
}
