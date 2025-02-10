package resource

import (
	"backend/shared/auth/domain/permission"
	"backend/shared/auth/domain/scope"
)

type Resources []Resource

func (rs Resources) Filter(filterAdmin bool) (resources []Resource) {
	for i := range rs {
		if rs[i].Name == defaultResource ||
			(filterAdmin && rs[i].Name == Admin) {
			continue
		}
		resources = append(resources, rs[i])
	}
	return
}

func (rs Resources) MapPermission(input permission.Permissions, scopes scope.Scopes) (rp []ResourcePermission) {
	for i := range rs {
		permsInResource := input.GetByResource(rs[i].ID)

		roles := make(map[string][]scope.Scope)
		for k := range permsInResource {
			for j := range permsInResource[k].Roles {
				roles[permsInResource[k].Roles[j].Name] = scopes.GetByName(permsInResource[k].Scopes)
			}
		}

		rpTemp := ResourcePermission{
			ID:          rs[i].ID,
			Name:        rs[i].Name,
			DisplayName: rs[i].DisplayName,
			Roles:       roles,
		}
		rp = append(rp, rpTemp)
	}
	return
}
