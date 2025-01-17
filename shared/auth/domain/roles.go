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

func DefaultRoles () []Role{
	return[]Role{
		{
		Name:        RoleAdmin,
		Description: "Administrator with tenant management capabilities",
	},

	{
		Name:        RoleWorker,
		Description: "Worker with basic access permissions",
	},	
	}

}

func (r Role) IsDefaultRole() bool {
	for _, defaultRole := range DefaultRoles() {
		if r.Name == defaultRole.Name {
			return true
		}
	}
	return false
}
