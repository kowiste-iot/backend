package role

import (
	"context"
	authCmd "ddd/shared/auth/domain/command"
	"ddd/shared/base/command"
	"ddd/shared/util"
)

type Role struct {
	ID          string
	Name        string
	Description string
}
type RoleProvider interface {
	// Role management methods
	CreateRole(ctx context.Context, input *authCmd.CreateRoleInput) (string, error)
	DeleteRole(ctx context.Context, input *authCmd.RoleIDInput) error
	GetRole(ctx context.Context, input *authCmd.RoleIDInput) (*Role, error)
	GetRoles(ctx context.Context, input *command.BaseInput) ([]Role, error)
	GetUserRoles(ctx context.Context, input *authCmd.UserRolesInput) ([]Role, error)

	// Role assignment methods
	AssignRoles(ctx context.Context, input *authCmd.AssignRolesInput) error
	RemoveRoles(ctx context.Context, input *authCmd.RemoveRolesInput) error
}

func NewRole(id string, Name string) *Role {
	return &Role{
		ID:   id,
		Name: Name,
	}
}

const (
	RoleAdmin = "admin"
	RoleUma   = "uma_protection"
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

func (r Role) IsDefaultRole() bool {
	for _, defaultRole := range AdminRoles() {
		if r.Name == defaultRole.Name {
			return true
		}
	}
	return false
}
func (r *Role) WithDescription(input *string) *Role {
	if input != nil {
		r.Description = *input
	}
	return r
}
