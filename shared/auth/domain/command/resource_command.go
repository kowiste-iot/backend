package command

import (
	"ddd/shared/auth/domain/scope"
	"ddd/shared/base/command"
	"fmt"
)

type UpdateResourceInput struct {
	command.BaseInput
	ID          string
	Name        string
	DisplayName string
	Roles       map[string][]scope.Scope
}

func ResourceName(roleName string) string {
	return fmt.Sprintf("%s-resource", roleName)
}
