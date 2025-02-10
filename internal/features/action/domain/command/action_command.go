package command

import "backend/shared/base/command"

type CreateActionInput struct {
	command.BaseInput
	ID          string `validate:"omitempty,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Parent      string `validate:"omitempty,uuidv7"`
	Enabled     bool
	Description string `validate:"omitempty,min=3,max=512"`
}

type UpdateActionInput struct {
	command.BaseInput
	ID          string `validate:"required,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Parent      string `validate:"omitempty,uuidv7"`
	Enabled     bool
	Description string `validate:"omitempty,min=3,max=512"`
}

type ActionIDInput struct {
	command.BaseInput
	ActionID string `validate:"required,uuidv7"`
}
