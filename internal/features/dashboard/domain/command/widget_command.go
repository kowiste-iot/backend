package command

import "backend/shared/base/command"

type CreateWidgetInput struct {
	command.BaseInput
	ID          string `validate:"omitempty,uuidv7"`
	DashboardID string `validate:"required,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
}

type UpdateWidgetInput struct {
	command.BaseInput
	ID          string `validate:"required,uuidv7"`
	DashboardID string `validate:"required,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
}

type WidgetIDInput struct {
	command.BaseInput
	DashboardID string `validate:"required,uuidv7"`
	WidgetID    string `validate:"required,uuidv7"`
}
