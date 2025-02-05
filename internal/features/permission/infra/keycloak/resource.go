package keycloak

import (
	"backend/internal/features/resource/domain"
	"backend/internal/features/resource/domain/command"
	baseCmd "backend/shared/base/command"
	"context"
	"fmt"

	"backend/shared/keycloak"

	"github.com/Nerzal/gocloak/v13"
)

type ResourceKeycloak struct {
	*keycloak.Keycloak
}

func New(core *keycloak.Keycloak) *ResourceKeycloak {
	return &ResourceKeycloak{
		Keycloak: core,
	}
}

func (rk ResourceKeycloak) CreateResource(ctx context.Context, input *baseCmd.BaseInput, resource domain.Resource) (*domain.Resource, error) {

	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	kcResource := gocloak.ResourceRepresentation{
		Name:        &resource.Name,
		DisplayName: &resource.DisplayName,
		Type:        &resource.Type,
		Scopes:      convertToScopeRepresentations(resource.Scopes),
	}

	createdResource, err := rk.Client.CreateResource(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		kcResource,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return convertFromKeycloakResource(createdResource), nil
}

func (rk ResourceKeycloak) GetResource(ctx context.Context, input *command.ResourceIDInput) (*domain.Resource, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}
	kcResource, err := rk.Client.GetResource(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		input.ResourceID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return convertFromKeycloakResource(kcResource), nil

}
func (rk ResourceKeycloak) ListResources(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Resource, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	kcResources, err := rk.Client.GetResources(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		gocloak.GetResourceParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	resources := make([]domain.Resource, len(kcResources))
	for i, kcResource := range kcResources {
		resources[i] = *convertFromKeycloakResource(kcResource)
	}

	return resources, nil
}

