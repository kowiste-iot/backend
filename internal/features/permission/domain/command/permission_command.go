package command

import (
	"backend/shared/base/command"
	"fmt"
)

type CreatePermissionInput struct {
	command.BaseInput
	ID               string
	Name             string
	Description      string
	Type             string
	ResourceType     string
	Resources        []string
	Scopes           []string
	Policies         []string
	DecisionStrategy string
}

func ResourceName(roleName string) string {
	return fmt.Sprintf("%s-resource", roleName)
}
