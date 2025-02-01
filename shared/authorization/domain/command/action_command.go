package command

import "ddd/shared/base/command"

type CreateActionInput struct {
	command.BaseInput
	Name        string
	DisplayName string
}

type UpdateActionInput struct {
	command.BaseInput
	Name        string
	DisplayName string
}
type ActionIDInput struct {
	command.BaseInput
	ActionID string
}
