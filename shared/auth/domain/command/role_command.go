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
	UserID string   `validate:"required,uuidv7"`
	Roles  []string `validate:"required,min=1,dive,required"`
}

type RemoveRolesInput struct {
	command.BaseInput
	UserID string   `validate:"required,uuidv7"`
	Roles  []string `validate:"required,min=1,dive,required"`
}

type UserRolesInput struct {
	command.BaseInput
	UserID string `validate:"required,uuidv7"`
}
