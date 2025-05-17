package domain

import (
	"backend/internal/features/user/domain"

	"errors"
)

type Permissions []Permission

func (rs Permissions) MapRoles(roles domain.Roles, filterAdmin bool) (permission []Permission, err error) {
	for i := range rs {
		if rs[i].Name == defaultPermission ||
			(filterAdmin && rs[i].Name == adminPermission) {
			continue
		}
		for j := range rs[i].Policies {
			role := roles.GetByName(domain.PolicyToRole(rs[i].Policies[j]))
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
		if rs[i].Resource == resourceID {
			p = append(p, rs[i])
		}
	}
	return
}
