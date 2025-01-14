package domain

import "context"

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
