package command

import (
	"ddd/shared/base/command"
	"fmt"
)

// type CreateRoleInput struct {
// 	command.BaseInput
// 	Name        string
// 	Description string
// }

//	type UpdateRoleInput struct {
//		command.BaseInput
//		Name        string
//		Description string
//	}
type PolicyNameInput struct {
	command.BaseInput
	PolicyName string
}

func PolicyName(roleName string) string {
	return fmt.Sprintf("%s-policy", roleName)
}
