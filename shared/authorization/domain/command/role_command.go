package command

import "ddd/shared/base/command"

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

