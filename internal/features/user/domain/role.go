package domain

import (
	"backend/internal/features/user/domain/command"
	baseCmd "backend/shared/base/command"
	"backend/shared/util"
	"context"
)

const (
	RoleAdmin = "admin"
	RoleUma   = "uma_protection" //TODO:move inside keycloak, never return
)

type RoleProvider interface {
	CreateRole(ctx context.Context, input *command.CreateRoleInput) (*Role, error)
	DeleteRole(ctx context.Context, input *command.RoleIDInput) error
	GetRole(ctx context.Context, input *command.RoleIDInput) (*Role, error)
	GetRoles(ctx context.Context, input *baseCmd.BaseInput) ([]Role, error)
}

type Role struct {
	ID          string
	PolicyID    string
	Name        string
	Description string
}

func NewRole(id string, Name string) *Role {
	return &Role{
		ID:   id,
		Name: Name,
	}
}

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

func (r Role) IsAdminRole() bool {
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
func (r *Role) WithPolicy(policyID string) *Role {
	r.PolicyID = policyID
	return r
}

type Roles []Role

func (r Roles) GetByName(roleName string) *Role {
	for i := range r {
		if r[i].Name == roleName {
			return &r[i]
		}
	}
	return nil
}
