package permission

import "ddd/shared/auth/domain/command"

type Permissions []Permission

func (rs Permissions) Filter(filterAdmin bool) (permission []Permission) {
	for i := range rs {
		if rs[i].Name == defaultPermission ||
			(filterAdmin && rs[i].Name == adminPermission) {
			continue
		}
		for j := range rs[i].Policies {
			rs[i].Roles = append(rs[i].Roles, command.PolicyToRole(rs[i].Policies[j]))
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
