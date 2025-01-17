package auth

type Role struct {
	ID          string
	Name        string
	Description string
}

const (
	RoleAdmin  = "admin"
	RoleWorker = "worker"
)

func NonAdminRoles() []Role {
	return []Role{
		{
			Name:        RoleWorker,
			Description: "Worker with basic access permissions",
		},
	}
}
func AdminRoles() []Role {
	return []Role{
		{
			Name:        RoleAdmin,
			Description: "Administrator with tenant management capabilities",
		},
	}
}

func AllRoles() []Role {
	return append(AdminRoles(), NonAdminRoles()...)

}

func (r Role) IsDefaultRole() bool {
	for _, defaultRole := range AllRoles() {
		if r.Name == defaultRole.Name {
			return true
		}
	}
	return false
}
