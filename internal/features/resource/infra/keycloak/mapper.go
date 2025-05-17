package keycloak

import (
	
	"backend/internal/features/resource/domain"

	"github.com/Nerzal/gocloak/v13"
)

func convertToScopeRepresentations(scopes []string) *[]gocloak.ScopeRepresentation {
	if len(scopes) == 0 {
		return nil
	}

	scopeReps := make([]gocloak.ScopeRepresentation, len(scopes))
	for i, scope := range scopes {
		name := scope
		scopeReps[i] = gocloak.ScopeRepresentation{
			Name: &name,
		}
	}
	return &scopeReps
}

func convertFromKeycloakResource(kr *gocloak.ResourceRepresentation) *domain.Resource {
	if kr == nil {
		return nil
	}

	resource := &domain.Resource{
		Name: *kr.Name,
	}

	// Safely handle optional fields
	if kr.ID != nil {
		resource.ID = *kr.ID
	}
	if kr.DisplayName != nil {
		resource.DisplayName = *kr.DisplayName
	}
	if kr.Type != nil {
		resource.Type = *kr.Type
	}

	if kr.Scopes != nil {
		scopes := make([]string, 0, len(*kr.Scopes))
		for _, scope := range *kr.Scopes {
			if scope.Name != nil {
				scopes = append(scopes, *scope.Name)
			}
		}
		resource.Scopes = scopes
	}

	return resource
}
