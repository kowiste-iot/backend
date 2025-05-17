package domain

import (
	"backend/internal/features/tenant/domain/command"
	userDomain "backend/internal/features/user/domain"
	baseCmd "backend/shared/base/command"
	"context"
)

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	Update(ctx context.Context, tenant *Tenant) error
	FindByID(ctx context.Context, tenantID string) (*Tenant, error)
	FindAll(ctx context.Context) ([]*Tenant, error)
	Remove(ctx context.Context, tenantID string) error
	HasChildren(ctx context.Context, tenantID string) (bool, error)
}
type BranchRepository interface {
	Create(ctx context.Context, tenantID string, branch *Branch) error
	Update(ctx context.Context, tenantID string, branch *Branch) error
	FindByID(ctx context.Context, tenantID, branchID string) (*Branch, error)
	FindAll(ctx context.Context, tenantID string) ([]*Branch, error)
	Remove(ctx context.Context, tenantID, branchID string) error
}
type TenantProvider interface {
	CreateTenant(ctx context.Context, tenant *Tenant) (string, error)
	UpdateTenant(ctx context.Context, tenant *Tenant) error
	DeleteTenant(ctx context.Context, tenantID string) error
	GetTenant(ctx context.Context, tenantID string) (*Tenant, error)
	CreateAdminGroup(ctx context.Context, tenantID string) (string, error)
}
type BranchProvider interface {
	CreateBranch(ctx context.Context, input *command.CreateBranchInput) (string, error)
	UpdateBranch(ctx context.Context, input *command.UpdateBranchInput) error
	DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error
	GetBranch(ctx context.Context, input *baseCmd.BaseInput) (*Branch, error)
	GetBranchUsers(ctx context.Context, input *baseCmd.BaseInput) ([]userDomain.User, error)
	AssignAdminsToBranch(ctx context.Context, input *baseCmd.BaseInput) error
	AssignUserToBranch(ctx context.Context, input *command.UserToBranch) error
	RemoveUserFromBranch(ctx context.Context, input *command.UserToBranch) error
}
