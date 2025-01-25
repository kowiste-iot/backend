package role

type Roles []Role

func (r Roles) GetByName(roleName string) *Role {
	for i := range r {
		if r[i].Name == roleName {
			return &r[i]
		}
	}
	return nil
}
