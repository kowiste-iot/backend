package app

import (
	"backend/internal/features/resource/domain"
	"backend/internal/features/resource/domain/command"

	appRole "backend/internal/features/role/app"
	rolesDomain "backend/internal/features/role/domain"
	scopeDomain "backend/internal/features/scope/domain"
	roleCmd "backend/internal/features/role/domain/command"

	appScope "backend/internal/features/scope/app"

	appPermission "backend/internal/features/permission/app"
	permissionDomain "backend/internal/features/permission/domain"


	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
)

type ResourceService interface {
	CreateResource(ctx context.Context, input *command.CreateResourceInput) (*domain.Resource, error)
	UpdateResource(ctx context.Context, input *command.UpdateResourceInput) (*domain.ResourcePermission, error)
	ListResources(ctx context.Context, input *baseCmd.BaseInput) ([]domain.ResourcePermission, error)
}
type Config struct {
	DefaultRoles []string
}

type ServiceDependencies struct {
	Repo       domain.ResourceProvider
	Roles      appRole.RoleService
	Scopes     appScope.ScopeService
	Permission appPermission.PermissionService
	Config     *Config
}
type resourceService struct {
	resourceProvider domain.ResourceProvider
	roles            appRole.RoleService
	scopes           appScope.ScopeService
	permission       appPermission.PermissionService
	config           *Config
	*base.BaseService
}

func NewService(base *base.BaseService, deps *ServiceDependencies) ResourceService {
	return &resourceService{
		resourceProvider: deps.Repo,
		roles:            deps.Roles,
		scopes:           deps.Scopes,
		permission:       deps.Permission,
		BaseService:      base,
		config:           deps.Config,
	}
}
func (s *resourceService) CreateResource(ctx context.Context, input *command.CreateResourceInput) (resource *domain.Resource, err error) {
	// err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: input.BaseInput,
	// 	Resource:  domain.ResourceR,
	// 	Scope:     scope.Create,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	r, err := domain.New(command.ResourceName(input.Name), input.Type, input.Scopes, input.DisplayName)
	if err != nil {
		return
	}
	resource, err = s.resourceProvider.CreateResource(ctx, &input.BaseInput, *r)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return resource, nil
}

func (s *resourceService) ListResources(ctx context.Context, input *baseCmd.BaseInput) (resources []domain.ResourcePermission, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  domain.ResourceR,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}

	var tempResources domain.Resources
	tempResources, err = s.resourceProvider.ListResources(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	tempResources = tempResources.Filter(true)
	var tempPermissions permissionDomain.Permissions

	tempPermissions, err = s.permission.ListPermissions(ctx, input)
	if err != nil {
		return
	}

	roles, err := s.roles.ListRoles(ctx, input)
	if err != nil {
		return nil, err
	}
	permissions, err := tempPermissions.MapRoles(roles, true)
	if err != nil {
		return
	}
	scopes, err := s.scopes.ListScopes(ctx, input)
	if err != nil {
		return
	}
	resources = tempResources.MapPermission(permissions, scopes)
	return resources, nil
}

func (s *resourceService) UpdateResource(ctx context.Context, input *command.UpdateResourceInput) (resource *domain.ResourcePermission, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  domain.ResourceR,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return nil, err
	}

	r, err := s.resourceProvider.GetResource(ctx, &command.ResourceIDInput{
		BaseInput:  input.BaseInput,
		ResourceID: input.ID,
	})
	if err != nil {
		return
	}
	mRoles := make(map[string]rolesDomain.Role)
	//get policty of role
	roles, err := s.roles.ListRoles(ctx, &input.BaseInput)
	if err != nil {
		return
	}
	for i := range roles {
		mRoles[roles[i].Name] = roles[i]
	}

	inputAssign := roleCmd.ResourceAssignRoleInput{
		BaseInput:    input.BaseInput,
		ResourceID:   r.ID,
		ResourceName: r.DisplayName,
	}
	err = s.roles.RemoveRolesFromResource(ctx, &inputAssign)
	if err != nil {
		return
	}

	//create permissions should be Assign Role to Resource ? maybe assign roles should be on permission
	for name, scopes := range input.Roles {
		role, ok := mRoles[name]
		if !ok {
			return nil, fmt.Errorf("error ")
		}
		inputAssign.RoleID = role.ID
		inputAssign.RoleName = name
		inputAssign.Scopes = scopes
		err = s.roles.AssignRoleToResource(ctx, &inputAssign)
		if err != nil {
			return
		}
	}

	return
}
