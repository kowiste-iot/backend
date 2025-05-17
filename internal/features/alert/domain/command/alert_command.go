package command

import "backend/shared/base/command"

type CreateAlertInput struct {
	command.BaseInput
	ID          string `validate:"omitempty,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Parent      string `validate:"omitempty,uuidv7"`
	Enabled     bool
	Description string `validate:"omitempty,min=3,max=512"`
}

type UpdateAlertInput struct {
	command.BaseInput
	ID          string `validate:"required,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Parent      string `validate:"omitempty,uuidv7"`
	Enabled     bool
	Description string `validate:"omitempty,min=3,max=512"`
}

type AlertIDInput struct {
	command.BaseInput
	AlertID string `validate:"required,uuidv7"`
}
