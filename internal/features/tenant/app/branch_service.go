package app

import (
	"backend/internal/features/tenant/domain"
	"backend/internal/features/tenant/domain/command"
	"backend/internal/features/user/dto"
	"backend/pkg/config"
	"backend/shared/auth/domain/role"
	"backend/shared/util"

	"backend/shared/base"
	baseCmd "backend/shared/base/command"

	appRole "backend/internal/features/role/app"
	roleCmd "backend/internal/features/role/domain/command"

	appResource "backend/internal/features/resource/app"
	resourceDomain "backend/internal/features/resource/domain"
	resourceCmd "backend/internal/features/resource/domain/command"

	appScope "backend/internal/features/scope/app"
	scopeDomain "backend/internal/features/scope/domain"
	scopeCmd "backend/internal/features/scope/domain/command"

	"backend/shared/validator"
	"context"
	"fmt"
)

type BranchService interface {
	CreateBranch(ctx context.Context, input *command.CreateBranchInput) (*domain.Branch, error)
	GetBranch(ctx context.Context, input *baseCmd.BaseInput) (*domain.Branch, error)
	ListBranches(ctx context.Context, tenantID string) ([]*domain.Branch, error)
	UpdateBranch(ctx context.Context, input *command.UpdateBranchInput) (*domain.Branch, error)
	DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error
	GetBranchUsers(ctx context.Context, input *baseCmd.BaseInput) ([]dto.UserDTO, error)
	AssignUserToBranch(ctx context.Context, input *command.UserToBranch) error
	RemoveUserFromBranch(ctx context.Context, input *command.UserToBranch) error
}
type BranchDependencies struct {
	Branch   domain.BranchProvider
	Role     appRole.RoleService
	Scope    appScope.ScopeService
	Resource appResource.ResourceService
	Repo     domain.BranchRepository
	Config   *config.TenantConfiguration
}
type branchService struct {
	repo     domain.BranchRepository
	branch   domain.BranchProvider
	role     appRole.RoleService
	scope    appScope.ScopeService
	resource appResource.ResourceService
	config   *config.TenantConfiguration
	*base.BaseService
}

func NewBranchService(base *base.BaseService, dep *BranchDependencies) BranchService {
	return &branchService{
		repo:        dep.Repo,
		branch:      dep.Branch,
		role:        dep.Role,
		scope:       dep.Scope,
		resource:    dep.Resource,
		config:      dep.Config,
		BaseService: base,
	}
}

func (s *branchService) CreateBranch(ctx context.Context, input *command.CreateBranchInput) (*domain.Branch, error) {
	// err := s.CheckPermission(ctx, resource.Tenant, scope.Create)
	// if err != nil {
	// 	return nil, err
	// }
	err := validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	branch, err := domain.NewBranch(input.TenantDomain, input.Name, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	// Create auth group
	branchID, err := s.branch.CreateBranch(ctx, &command.CreateBranchInput{
		TenantDomain: input.TenantDomain,
		Name:         input.Name,
		Description:  input.Description,
		Default:      input.Default,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auth branch: %w", err)
	}
	baseInput := baseCmd.BaseInput{
		TenantDomain: input.TenantDomain,
		BranchName:   input.Name,
	}
	branch.SetAuthBranchID(branchID)
	//Assign admins to created branch
	err = s.branch.AssignAdminsToBranch(ctx, &baseInput)
	if err != nil {
		return nil, err
	}
	//add scopes to group
	for _, scName := range scopeDomain.AllScopes() {
		_, err := s.scope.CreateScope(ctx, &scopeCmd.CreateScopeInput{
			BaseInput:   baseInput,
			ID:          branchID,
			Name:        scName,
			DisplayName: util.CapitalizeFirst(scName),
		})
		if err != nil {
			return nil, err
		}
	}

	//Roles, create default roles
	for _, role := range role.AllRoles(s.config.Authorization.Roles) {
		input := roleCmd.CreateRoleInput{
			BaseInput:   baseInput,
			Name:        role.Name,
			Description: role.Description,
		}
		err = s.role.CreateDefaultRoles(ctx, &input)
		if err != nil {
			return nil, fmt.Errorf("failed to create auth role: %w", err)
		}
	}

	//create resources (and permissions) in branch
	resources := resourceDomain.EndpointsResources(s.config.Authorization.Resources)
	for _, res := range resources {
		_, err = s.resource.CreateResource(ctx, &resourceCmd.CreateResourceInput{
			BaseInput:   baseInput,
			Name:        resourceCmd.ResourceName(res.Name),
			DisplayName: res.Name,
			Type:        res.Type,
			Scopes:      res.Scopes,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := s.repo.Create(ctx, input.TenantDomain, branch); err != nil {
		// Cleanup auth group if DB save fails
		i := baseCmd.NewInput(input.TenantDomain, branchID)
		if delErr := s.branch.DeleteBranch(ctx, &i); delErr != nil {
			s.Logger.Error(ctx, err, "failed to cleanup auth group after branch creation failure", nil)
		}
		return nil, fmt.Errorf("failed to save branch: %w", err)
	}

	return branch, nil
}

func (s *branchService) GetBranch(ctx context.Context, input *baseCmd.BaseInput) (*domain.Branch, error) {

	branch, err := s.repo.FindByID(ctx, input.TenantDomain, input.BranchName)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}
	return branch, nil
}

func (s *branchService) ListBranches(ctx context.Context, tenantID string) ([]*domain.Branch, error) {

	branches, err := s.repo.FindAll(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}
	return branches, nil
}

func (s *branchService) UpdateBranch(ctx context.Context, cmd *command.UpdateBranchInput) (*domain.Branch, error) {

	branch, err := s.repo.FindByID(ctx, cmd.TenantDomain, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	err = branch.Update(cmd.Name, cmd.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to update branch: %w", err)
	}

	// Update Keycloak group
	err = s.branch.UpdateBranch(ctx, &command.UpdateBranchInput{
		TenantDomain: cmd.TenantDomain,
		ID:           cmd.ID,
		Name:         cmd.Name,
		Description:  cmd.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update auth group: %w", err)
	}

	if err := s.repo.Update(ctx, cmd.TenantDomain, branch); err != nil {
		return nil, fmt.Errorf("failed to save branch: %w", err)
	}

	return branch, nil
}

func (s *branchService) DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error {

	_, err := s.repo.FindByID(ctx, input.TenantDomain, input.BranchName)
	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}

	err = s.branch.DeleteBranch(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete auth branch: %w", err)
	}

	if err := s.repo.Remove(ctx, input.TenantDomain, input.BranchName); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}
func (s *branchService) GetBranchUsers(ctx context.Context, input *baseCmd.BaseInput) ([]dto.UserDTO, error) {
	return nil, nil
}
func (s *branchService) AssignUserToBranch(ctx context.Context, input *command.UserToBranch) error {
	return nil
}
func (s *branchService) RemoveUserFromBranch(ctx context.Context, input *command.UserToBranch) error {
	return nil
}
