package keycloak

import (
	"context"
	"ddd/shared/auth/domain/resource"
	baseCmd "ddd/shared/base/command"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

// ResourceProvider interface methods for KeycloakService
func (ks *KeycloakService) CreateResource(ctx context.Context, tenantID, clientID string, resource resource.Resource) (*resource.Resource, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	kcResource := gocloak.ResourceRepresentation{
		Name:        &resource.Name,
		DisplayName: &resource.DisplayName,
		Type:        &resource.Type,
		URIs:        &resource.URIs,
		Scopes:      convertToScopeRepresentations(resource.Scopes),
		Attributes:  &resource.Attributes,
		IconURI:     &resource.IconURI,
	}

	createdResource, err := ks.client.CreateResource(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		kcResource,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return convertFromKeycloakResource(createdResource), nil
}

func (ks *KeycloakService) UpdateResource(ctx context.Context, tenantID, clientID string, resource resource.Resource) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	kcResource := gocloak.ResourceRepresentation{
		ID:          &resource.ID,
		Name:        &resource.Name,
		DisplayName: &resource.DisplayName,
		Type:        &resource.Type,
		URIs:        &resource.URIs,
		Scopes:      convertToScopeRepresentations(resource.Scopes),
		Attributes:  &resource.Attributes,
		IconURI:     &resource.IconURI,
	}

	err = ks.client.UpdateResource(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		kcResource,
	)
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	return nil
}

func (ks *KeycloakService) DeleteResource(ctx context.Context, tenantID, clientID, resourceID string) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.DeleteResource(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		resourceID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return nil
}

func (ks *KeycloakService) GetResource(ctx context.Context, tenantID, clientID, resourceID string) (*resource.Resource, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	kcResource, err := ks.client.GetResource(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		resourceID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return convertFromKeycloakResource(kcResource), nil
}

func (ks *KeycloakService) ListResources(ctx context.Context, input *baseCmd.BaseInput) ([]resource.Resource, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.fetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	kcResources, err := ks.client.GetResources(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		gocloak.GetResourceParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	resources := make([]resource.Resource, len(kcResources))
	for i, kcResource := range kcResources {
		resources[i] = *convertFromKeycloakResource(kcResource)
	}

	return resources, nil
}

// Helper functions
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

func convertFromKeycloakResource(kr *gocloak.ResourceRepresentation) *resource.Resource {
	if kr == nil {
		return nil
	}

	resource := &resource.Resource{
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
	if kr.IconURI != nil {
		resource.IconURI = *kr.IconURI
	}
	if kr.URIs != nil {
		resource.URIs = *kr.URIs
	}
	if kr.Attributes != nil {
		resource.Attributes = *kr.Attributes
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
