package app

import (
	"backend/internal/features/tenant/domain"
	"backend/internal/features/tenant/domain/command"

	appUser "backend/internal/features/user/app"
	userCmd "backend/internal/features/user/domain/command"
	roleDomain "backend/internal/features/role/domain"

	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"

	"errors"
	"fmt"
)

type TenantService interface {
	CreateTenant(ctx context.Context, cmd *command.CreateTenantInput) (*domain.Tenant, error)
	GetTenant(ctx context.Context, tenantID string) (*domain.Tenant, error)
	ListTenants(ctx context.Context) ([]*domain.Tenant, error)
	UpdateTenant(ctx context.Context, cmd *command.UpdateTenantInput) (*domain.Tenant, error)
	DeleteTenant(ctx context.Context, tenantID string) error
}
type ServiceDependencies struct {
	Branch BranchService
	User   appUser.UserService
	Tenant   domain.TenantProvider
	Repo   domain.TenantRepository
}

type tenantService struct {
	repo   domain.TenantRepository
	tenant   domain.TenantProvider
	user   appUser.UserService
	branch BranchService
	*base.BaseService
}

func NewTenantService(base *base.BaseService, dep *ServiceDependencies) TenantService {
	return &tenantService{
		repo:        dep.Repo,
		tenant:        dep.Tenant,
		user:        dep.User,
		branch:      dep.Branch,
		BaseService: base,
	}
}
func (s tenantService) CreateTenant(ctx context.Context, input *command.CreateTenantInput) (*domain.Tenant, error) {
	// err := s.CheckPermission(ctx, resource.Tenant, scope.Create)
	// if err != nil {
	// 	return nil, err
	// }
	err := validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	tenant, err := domain.NewTenant(input.Name, input.Domain, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	//create auth teanat
	id, err := s.tenant.CreateTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}
	tenant.SetAuthID(id)

	//Create default branch
	defaultB := command.CreateBranchInput{
		TenantDomain: tenant.Domain(),
		Name:         input.Branch,
		Description:  "Default Branch",
	}
	createdBranch, err := s.branch.CreateBranch(ctx, &defaultB)
	if err != nil {
		return nil, fmt.Errorf("failed to create default branch: %w", err)
	}

	//Create admins branch they are tenant admin
	adminB := command.CreateBranchInput{
		TenantDomain: tenant.Domain(),
		Name:         domain.AdminBranch,
		Description:  "Admin Group",
	}
	adminBranch, err := s.branch.CreateBranch(ctx, &adminB)
	if err != nil {
		return nil, fmt.Errorf("failed to create default branch: %w", err)
	}

	// Create Admin user
	u := userCmd.CreateUserInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), input.Branch),
		Email:     input.AdminEmail,
		FirstName: "admin",
		LastName:  "user",
		Roles:     []string{roleDomain.RoleAdmin},
	}
	user, err := s.user.CreateUser(ctx, &u)
	if err != nil {
		return nil, fmt.Errorf("failed to create default branch: %w", err)
	}

	//Assign to branch
	ub := command.UserToBranch{
		TenantDomain: tenant.Domain(),
		UserID:       user.AuthID,
		Branchs:      []string{createdBranch.AuthBranchID(), adminBranch.AuthBranchID()},
	}
	err = s.branch.AssignUserToBranch(ctx, &ub)
	if err != nil {
		return nil, fmt.Errorf("failed to assign admin to branch: %w", err)
	}
	//save in repo
	if err := s.repo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

func (s *tenantService) GetTenant(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	// err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{} resource.Tenant, scope.View)
	// if err != nil {
	// 	return nil, err
	// }
	tenant, err := s.repo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	return tenant, nil
}

func (s *tenantService) ListTenants(ctx context.Context) ([]*domain.Tenant, error) {
	// err := s.CheckPermission(ctx, resource.Tenant, scope.View)
	// if err != nil {
	// 	return nil, err
	// }
	tenants, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return tenants, nil
}

func (s *tenantService) UpdateTenant(ctx context.Context, cmd *command.UpdateTenantInput) (*domain.Tenant, error) {
	// err := s.CheckPermission(ctx, resource.Tenant, scope.Update)
	// if err != nil {
	// 	return nil, err
	// }
	tenant, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	err = tenant.Update(cmd.Name, cmd.Domain, cmd.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	if err := s.repo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return tenant, nil
}
func (s *tenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	// err := s.CheckPermission(ctx, resource.Tenant, scope.Delete)
	// if err != nil {
	// 	return err
	// }
	has, err := s.repo.HasChildren(ctx, tenantID)
	if err != nil {
		return err
	}
	if has {
		return errors.New("tenant has children and cannot be deleted")
	}
	err = s.repo.Remove(ctx, tenantID)
	if err != nil {
		return err
	}
	return nil
}
