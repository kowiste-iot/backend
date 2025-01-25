package app

import (
	"context"
	"ddd/shared/auth/domain/command"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/resource"
	baseCmd "ddd/shared/base/command"
)

func (s *Service) GetResources(ctx context.Context, input *baseCmd.BaseInput) (resources []resource.ResourcePermission, err error) {

	var tempResources resource.Resources
	tempResources, err = s.resourceProvider.ListResources(ctx, input)
	if err != nil {
		return
	}

	tempResources = tempResources.Filter(true)
	var tempPermissions permission.Permissions

	tempPermissions, err = s.permissionProvider.ListPermissions(ctx, input)
	if err != nil {
		return
	}
	roles, err := s.tenantProvider.GetRoles(ctx, input)
	if err != nil {
		return
	}
	permissions, err := tempPermissions.MapRoles(roles, true)
	if err != nil {
		return
	}
	scopes, err := s.scopeProvider.ListScopes(ctx, input)
	if err != nil {
		return
	}
	resources = tempResources.MapPermission(permissions, scopes)
	return
}

func (s *Service) UpdateResource(ctx context.Context, input *command.UpdateResourceInput) (resurce resource.ResourcePermission, err error) {
	
	//get policty of role
	//get resource
	//delete permissions
	//create permissions

	return
}
