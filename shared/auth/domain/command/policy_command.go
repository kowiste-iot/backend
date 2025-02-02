package command

import (
	"backend/shared/base/command"
	"fmt"
	"strings"
)

const (
	defaultName string = "-policy"
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
	return fmt.Sprintf("%s"+defaultName, roleName)
}
func PolicyToRole(policyName string) string {
	return strings.Replace(policyName, defaultName, "", -1)
}
