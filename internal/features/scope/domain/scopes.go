package domain

type Scopes []Scope

func (s Scopes) GetByName(scopesName []string) []Scope {
	out := make([]Scope, 0)
	for i := range s {
		for j := range scopesName {
			if s[i].Name == scopesName[j] {
				out = append(out, s[i])
			}
		}

	}
	return out
}
