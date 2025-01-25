package permission

import (
	"ddd/shared/auth/domain/command"
	"ddd/shared/auth/domain/role"
	"errors"
)

type Permissions []Permission

func (rs Permissions) MapRoles(roles role.Roles, filterAdmin bool) (permission []Permission, err error) {
	for i := range rs {
		if rs[i].Name == defaultPermission ||
			(filterAdmin && rs[i].Name == adminPermission) {
			continue
		}
		for j := range rs[i].Policies {
			role := roles.GetByName(command.PolicyToRole(rs[i].Policies[j]))
			if role == nil {
				return nil, errors.New("role not found")
			}
			rs[i].Roles = append(rs[i].Roles, *role)
		}
		permission = append(permission, rs[i])
	}
	return
}

func (rs Permissions) GetByResource(resourceID string) (p []Permission) {

	for i := range rs {
		for j := range rs[i].Resources {
			if rs[i].Resources[j] == resourceID {
				p = append(p, rs[i])
			}
		}
	}
	return
}
