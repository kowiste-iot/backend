package auth

import "ddd/shared/util"

type Role struct {
	ID          string
	Name        string
	Description string
}

const (
	RoleAdmin = "admin"
)

func NonAdminRoles(input []string) (roles []Role) {
	for i := range input {
		roles = append(roles, Role{
			Name:        input[i],
			Description: util.CapitalizeFirst(input[i]) + " with basic access permissions",
		})
	}
	return
}
func AdminRoles() []Role {
	return []Role{
		{
			Name:        RoleAdmin,
			Description: "Administrator with tenant management capabilities",
		},
	}
}

func AllRoles(nonAdminRoles []string) []Role {
	return append(AdminRoles(), NonAdminRoles(nonAdminRoles)...)

}

func (r Role) IsDefaultRole(nonAdminRoles []string) bool {
	for _, defaultRole := range AllRoles(nonAdminRoles) {
		if r.Name == defaultRole.Name {
			return true
		}
	}
	return false
}
