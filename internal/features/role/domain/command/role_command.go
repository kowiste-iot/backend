package command

import "backend/shared/base/command"

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
