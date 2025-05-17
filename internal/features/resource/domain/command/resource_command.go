package command

import (
	scopeDomain "backend/internal/features/scope/domain"

	"backend/shared/base/command"
	"fmt"
)

type ResourceIDInput struct {
	command.BaseInput
	ResourceID string `validate:"required,uuid"`
}

type UpdateResourceInput struct {
	command.BaseInput
	ID          string
	Name        string
	DisplayName string
	Roles       map[string][]scopeDomain.Scope
}
type CreateResourceInput struct {
	command.BaseInput
	ID          string
	Name        string
	Type        string
	DisplayName string
	Scopes      []string
}

func ResourceName(roleName string) string {
	return fmt.Sprintf("%s-resource", roleName)
}
