package command

import "fmt"

func ResourceName(roleName string) string {
	return fmt.Sprintf("%s-resource", roleName)
}
