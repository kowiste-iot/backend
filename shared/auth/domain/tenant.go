package auth

import (
	"context"
	authCmd "ddd/shared/auth/domain/command"
	"ddd/shared/base/command"
)

type TenantProvider interface {
	CreateTenant(ctx context.Context, tenant *Tenant) (string, error)
	UpdateTenant(ctx context.Context, tenant *Tenant) error
	DeleteTenant(ctx context.Context, tenantID string) error
	GetTenant(ctx context.Context, tenantID string) (*Tenant, error)
	// Branch management methods
	CreateBranch(ctx context.Context, input *authCmd.CreateBranchInput) (string, error)
	UpdateBranch(ctx context.Context, input *authCmd.UpdateBranchInput) error
	DeleteBranch(ctx context.Context, input *command.BaseInput) error
	GetBranch(ctx context.Context, input *command.BaseInput) (*Branch, error)
	GetBranchUsers(ctx context.Context, input *command.BaseInput) ([]User, error)
	AssignUserToBranch(ctx context.Context, input *authCmd.UserToBranch) error
	RemoveUserFromBranch(ctx context.Context, input *authCmd.UserToBranch) error

	// Role management methods
	CreateRole(ctx context.Context, input *authCmd.CreateRoleInput) (string, error)
	DeleteRole(ctx context.Context, input *authCmd.RoleIDInput) error
	GetRole(ctx context.Context, input *authCmd.RoleIDInput) (*Role, error)
	GetRoles(ctx context.Context, input *command.BaseInput) ([]Role, error)
	// Role assignment methods
	AssignRoles(ctx context.Context, input *authCmd.AssignRolesInput) error
	RemoveRoles(ctx context.Context, input *authCmd.RemoveRolesInput) error
	GetUserRoles(ctx context.Context, input *authCmd.UserRolesInput) ([]Role, error)
}

type Tenant struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Domain   string         `json:"domain"`
	Settings TenantSettings `json:"settings"`
	RealmID  *string        `json:"realmId,omitempty"`
	Theme    *TenantTheme   `json:"tenatTheme,omitempty"`
}
type TenantTheme struct {
	Login string `json:"login,omitempty"`
}

type TenantSettings struct {
	AllowedDomains []string          `json:"allowedDomains"`
	Features       map[string]bool   `json:"features"`
	CustomConfig   map[string]string `json:"customConfig"`
}

