package app

import (
	"context"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/resource"
	baseCmd "ddd/shared/base/command"
	"fmt"
)

func (s *Service) GetResources(ctx context.Context, input *baseCmd.BaseInput) (resources []resource.Resource, err error) {

	var tempResources resource.Resources
	tempResources, err = s.resourceProvider.ListResources(ctx, input)
	if err != nil {
		return
	}

	resources = tempResources.Filter(true)
	var tempPermissions permission.Permissions

	tempPermissions, err = s.permissionProvider.ListPermissions(ctx, input)
	if err != nil {
		return
	}
	permissions:=tempPermissions.Filter(true)
	fmt.Println("re", permissions)
	return 
}
