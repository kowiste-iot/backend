package app

import (
	"context"
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
	permissions := tempPermissions.Filter(true)
	resources = tempResources.MapPermission(permissions)
	return
}
