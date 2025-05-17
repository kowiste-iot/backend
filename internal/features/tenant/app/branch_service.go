package app

import (
	"backend/internal/features/tenant/domain"
	"backend/internal/features/tenant/domain/command"
	"backend/internal/features/user/dto"
	"backend/pkg/config"

	"backend/shared/util"

	"backend/shared/base"
	baseCmd "backend/shared/base/command"

	appRole "backend/internal/features/user/app"
	roleDomain "backend/internal/features/user/domain"
	roleCmd "backend/internal/features/user/domain/command"

	appResource "backend/internal/features/resource/app"
	resourceDomain "backend/internal/features/resource/domain"
	resourceCmd "backend/internal/features/resource/domain/command"

	appPermission "backend/internal/features/permission/app"
	permissionDomain "backend/internal/features/permission/domain"
	permissionCmd "backend/internal/features/permission/domain/command"

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
	Branch     domain.BranchProvider
	Role       appRole.RoleService
	Scope      appScope.ScopeService
	Resource   appResource.ResourceService
	Permission appPermission.PermissionService
	Repo       domain.BranchRepository
	Config     *config.TenantConfiguration
}
type branchService struct {
	repo       domain.BranchRepository
	branch     domain.BranchProvider
	role       appRole.RoleService
	scope      appScope.ScopeService
	resource   appResource.ResourceService
	permission appPermission.PermissionService
	config     *config.TenantConfiguration
	*base.BaseService
}

func NewBranchService(base *base.BaseService, dep *BranchDependencies) BranchService {
	return &branchService{
		repo:        dep.Repo,
		branch:      dep.Branch,
		role:        dep.Role,
		scope:       dep.Scope,
		resource:    dep.Resource,
		permission:  dep.Permission,
		config:      dep.Config,
		BaseService: base,
	}
}

func (s *branchService) CreateBranch(ctx context.Context, input *command.CreateBranchInput) (branch *domain.Branch, err error) {
	// err := s.CheckPermission(ctx, resource.Tenant, scope.Create)
	// if err != nil {
	// 	return nil, err
	// }
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	branch, err = domain.NewBranch(input.TenantDomain, input.Name, input.Description)
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

	if !input.Default {
		//Assign admins to created branch when is not default branch
		err = s.branch.AssignAdminsToBranch(ctx, &baseInput)
		if err != nil {
			return nil, err
		}
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
	policies := make(map[string]string)

	//Roles, create default roles
	for _, role := range roleDomain.AllRoles(s.config.Authorization.Roles) {
		input := roleCmd.CreateRoleInput{
			BaseInput:   baseInput,
			Name:        role.Name,
			Description: role.Description,
		}
		r, err := s.role.CreateDefaultRoles(ctx, &input)
		if err != nil {
			return nil, fmt.Errorf("failed to create auth role: %w", err)
		}
		policies[r.Name] = r.PolicyID
	}

	//create resources (and permissions) in branch
	resources := resourceDomain.EndpointsResources(s.config.Authorization.Resources)
	for _, res := range resources {
		createdResource, err := s.resource.CreateResource(ctx, &resourceCmd.CreateResourceInput{
			BaseInput:   baseInput,
			Name:        resourceCmd.ResourceName(res.Name),
			DisplayName: res.Name,
			Type:        res.Type,
			Scopes:      res.Scopes,
		})
		if err != nil {
			return nil, err
		}
		resourceConfig, exist := s.config.Authorization.Resources[resourceCmd.ResourceName(res.Name)]
		if !exist {
			return nil, fmt.Errorf("error fetching resource %s", res.Name)
		}
		for roleName, scopes := range resourceConfig.Permissions {
			id, exist := policies[roleName]
			if !exist {
				return nil, fmt.Errorf("error fetching policy for  %s", roleName)
			}
			//create permissions
			_, err = s.permission.CreatePermission(ctx, &permissionCmd.CreatePermissionInput{
				BaseInput:        baseInput,
				Name:             permissionDomain.NameNonAdmin(roleName, res.Name),
				Description:      fmt.Sprintf("Permission for %s resource with %s role", res.Name, roleName),
				Type:             permissionDomain.TypeScope,
				Resources:        createdResource.ID,
				Scopes:           scopes,
				Policies:         []string{id},
				DecisionStrategy: permissionDomain.DecisionAffirmative,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	//resource permission for admin
	sc := scopeDomain.AllScopes()

	pID := policies[roleDomain.RoleAdmin]
	perm := permissionCmd.CreatePermissionInput{
		BaseInput:        baseInput,
		Name:             permissionDomain.NameAdmin(),
		Description:      fmt.Sprintf("Permission for %s resource with %s role", roleDomain.RoleAdmin, roleDomain.RoleAdmin),
		ResourceType:     s.config.Authorization.AdminGroup,
		Type:             permissionDomain.TypeResource,
		Scopes:           sc,
		Policies:         []string{pID},
		DecisionStrategy: permissionDomain.DecisionAffirmative,
	}

	_, err = s.permission.CreatePermission(ctx, &perm)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission for %s: %w", roleDomain.RoleAdmin, err)
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

	return s.branch.AssignUserToBranch(ctx, input)
}
func (s *branchService) RemoveUserFromBranch(ctx context.Context, input *command.UserToBranch) error {
	return s.branch.RemoveUserFromBranch(ctx, input)
}
