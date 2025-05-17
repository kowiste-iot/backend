package command

import (
	"backend/shared/base/command"
	"fmt"
)

type CreateScopeInput struct {
	command.BaseInput
	ID          string
	Name        string
	DisplayName string
}

type ScopeIDInput struct {
	command.BaseInput
	ResourceID string `validate:"required,uuid"`
}

func ResourceName(roleName string) string {
	return fmt.Sprintf("%s-resource", roleName)
}
