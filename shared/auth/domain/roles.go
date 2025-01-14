package auth

type Role struct {
	Name        string
	Description string
}

const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleSupervisor = "supervisor"
	RoleWorker     = "worker"
)

var DefaultRoles = []Role{
	{
		Name:        RoleSuperAdmin,
		Description: "Super Administrator with full system access",
	},
	{
		Name:        RoleAdmin,
		Description: "Administrator with tenant management capabilities",
	},
	{
		Name:        RoleSupervisor,
		Description: "Supervisor with team management capabilities",
	},
	{
		Name:        RoleWorker,
		Description: "Worker with basic access permissions",
	},
}

func (r Role) IsDefaultRole() bool {
	for _, defaultRole := range DefaultRoles {
		if r.Name == defaultRole.Name {
			return true
		}
	}
	return false
}
