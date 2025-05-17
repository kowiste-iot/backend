package command

import "backend/shared/base/command"

type CreateUserInput struct {
	command.BaseInput
	ID        string   `validate:"required,uuidv7"`
	Email     string   `validate:"required,email"`
	FirstName string   `validate:"required,min=3,max=255"`
	LastName  string   `validate:"required,min=3,max=255"`
	Roles     []string `validate:"required,roles"`
}

type UpdateUserInput struct {
	command.BaseInput
	ID        string `validate:"required,uuidv7"`
	Email     string `validate:"required,email"`
	FirstName string `validate:"required,min=3,max=255"`
	LastName  string `validate:"required,min=3,max=255"`
	Roles     []string
}
type UserIDInput struct {
	command.BaseInput
	UserID string `validate:"required,uuidv7"`
}

type UserRolesInput struct {
	command.BaseInput
	UserID string
}
