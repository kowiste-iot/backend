package app

import (
	"context"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/resource"
	resourceCmd "ddd/shared/auth/domain/resource/command"
	"ddd/shared/auth/domain/role"
	baseCmd "ddd/shared/base/command"
	"fmt"
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
	roles, err := s.roleProvider.GetRoles(ctx, input)
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

func (s *Service) UpdateResource(ctx context.Context, input *resourceCmd.UpdateResourceInput) (res *resource.ResourcePermission, err error) {

	r, err := s.resourceProvider.GetResource(ctx, &resourceCmd.ResourceIDInput{
		BaseInput:  input.BaseInput,
		ResourceID: input.ID,
	})
	if err != nil {
		return
	}
	mRoles := make(map[string]role.Role)
	//get policty of role
	roles, err := s.GetRoles(ctx, &input.BaseInput)
	if err != nil {
		return
	}
	for i := range roles {
		mRoles[roles[i].Name] = roles[i]
	}

	inputAssign := resourceCmd.ResourceAssignRoleInput{
		BaseInput:    input.BaseInput,
		ResourceID:   r.ID,
		ResourceName: r.DisplayName,
	}
	err = s.resourceProvider.RemoveRolesFromResource(ctx, &inputAssign)
	if err != nil {
		return
	}

	//create permissions shoudl be Assign Role to Resource
	for name, scopes := range input.Roles {
		role, ok := mRoles[name]
		if !ok {
			return nil, fmt.Errorf("error ")
		}
		inputAssign.RoleID = role.ID
		inputAssign.RoleName = name
		inputAssign.Scopes = scopes
		err = s.resourceProvider.AssignRoleToResource(ctx, &inputAssign)
		if err != nil {
			return
		}
	}
	//TODO: resource should return a resource permission?
	res = new(resource.ResourcePermission)
	return
}
