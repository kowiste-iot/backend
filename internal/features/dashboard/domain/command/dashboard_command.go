package command

import "backend/shared/base/command"

type CreateDashboardInput struct {
	command.BaseInput
	ID          string `validate:"omitempty,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Parent      string `validate:"omitempty,uuidv7"`
	Description string `validate:"omitempty,min=3,max=512"`
}

type UpdateDashboardInput struct {
	command.BaseInput
	ID          string `validate:"required,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Parent      string `validate:"omitempty,uuidv7"`
	Description string `validate:"omitempty,min=3,max=512"`
}

type DashboardIDInput struct {
	command.BaseInput
	DashboardID string `validate:"required,uuidv7"`
}
