package keycloak

import (
	"backend/internal/features/permission/domain"
	roleDomain "backend/internal/features/user/domain"
)

type permissionKc struct {
	ID               string            `json:"id,omitempty"`
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	Type             string            `json:"type"` // resource-based, scope-based
	ResourceType     string            `json:"resourceType"`
	Resources        []string          `json:"resources,omitempty"`
	Scopes           []string          `json:"scopes,omitempty"`
	Policies         []string          `json:"policies"`
	Roles            []roleDomain.Role `json:"roles,omitempty"`
	DecisionStrategy string            `json:"decisionStrategy"` // UNANIMOUS, AFFIRMATIVE, CONSENSUS
	Logic            string            `json:"logic"`
}

func newPermissionKc(input *domain.Permission) *permissionKc {
	if input == nil {
		return nil
	}
	internal:= &permissionKc{
		ID:               input.ID,
		Name:             input.Name,
		Description:      input.Description,
		Type:             input.Type,
		ResourceType:     input.ResourceType,
		Scopes:           input.Scopes,
		Policies:         input.Policies,
		Roles:            input.Roles,
		DecisionStrategy: input.DecisionStrategy,
		Logic:            input.Logic,
	}
	if input.Resource!=""{
		internal.Resources=[]string{input.Resource}
	}
	return internal
}
func (p *permissionKc) ToDomain() *domain.Permission {
    if p == nil {
        return nil
    }

    var resource string
    if len(p.Resources) > 0 {
        resource = p.Resources[0]
    }

    return &domain.Permission{
        ID:               p.ID,
        Name:             p.Name,
        Description:      p.Description,
        Type:             p.Type,
        ResourceType:     p.ResourceType,
        Resource:         resource,
        Scopes:           p.Scopes,
        Policies:         p.Policies,
        Roles:            p.Roles,
        DecisionStrategy: p.DecisionStrategy,
        Logic:           p.Logic,
    }
}
// func convertToScopeRepresentations(scopes []string) *[]gocloak.ScopeRepresentation {
// 	if len(scopes) == 0 {
// 		return nil
// 	}

// 	scopeReps := make([]gocloak.ScopeRepresentation, len(scopes))
// 	for i, scope := range scopes {
// 		name := scope
// 		scopeReps[i] = gocloak.ScopeRepresentation{
// 			Name: &name,
// 		}
// 	}
// 	return &scopeReps
// }

// func convertFromKeycloakResource(kr *gocloak.ResourceRepresentation) *domain.Resource {
// 	if kr == nil {
// 		return nil
// 	}

// 	resource := &domain.Resource{
// 		Name: *kr.Name,
// 	}

// 	// Safely handle optional fields
// 	if kr.ID != nil {
// 		resource.ID = *kr.ID
// 	}
// 	if kr.DisplayName != nil {
// 		resource.DisplayName = *kr.DisplayName
// 	}
// 	if kr.Type != nil {
// 		resource.Type = *kr.Type
// 	}

// 	if kr.Scopes != nil {
// 		scopes := make([]string, 0, len(*kr.Scopes))
// 		for _, scope := range *kr.Scopes {
// 			if scope.Name != nil {
// 				scopes = append(scopes, *scope.Name)
// 			}
// 		}
// 		resource.Scopes = scopes
// 	}

// 	return resource
// }
