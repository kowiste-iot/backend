package command

import (
	scopeDomain "backend/internal/features/scope/domain"

	"backend/shared/base/command"
)

type CreateRoleInput struct {
	command.BaseInput
	Name        string
	Description string
}

type UpdateRoleInput struct {
	command.BaseInput
	Name        string
	Description string
}
type RoleIDInput struct {
	command.BaseInput
	RoleID string
}

type AssignRolesInput struct {
	command.BaseInput
	UserID string
	Roles  []string
}

type RemoveRolesInput struct {
	command.BaseInput
	UserID string
	Roles  []string
}
type ResourceAssignRoleInput struct {
	command.BaseInput
	ResourceID   string `validate:"required,uuid"`
	ResourceName string `validate:"uuid"`
	RoleID       string `validate:"required,uuid"`
	RoleName     string `validate:"uuid"`
	Scopes       []scopeDomain.Scope
}
