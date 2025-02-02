package keycloak

import (
	"backend/shared/auth/domain/command"
	"backend/shared/auth/domain/permission"
	"backend/shared/auth/domain/resource"
	resourceCmd "backend/shared/auth/domain/resource/command"
	baseCmd "backend/shared/base/command"
	"context"
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

func (ks *KeycloakService) GetResource(ctx context.Context, input *resourceCmd.ResourceIDInput) (*resource.Resource, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.fetchClient(ctx, &input.BaseInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}
	kcResource, err := ks.client.GetResource(
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

// roles asssign, unassign
func (ks *KeycloakService) AssignRoleToResource(ctx context.Context, input *resourceCmd.ResourceAssignRoleInput) (err error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.fetchClient(ctx, &input.BaseInput)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	fmt.Println(token)
	p, err := ks.GetPolicyByName(ctx, &command.PolicyNameInput{
		BaseInput:  input.BaseInput,
		PolicyName: command.PolicyName(input.RoleName),
	})
	if err != nil {
		return
	}
	scopes := make([]string, 0)
	for i := range input.Scopes {
		scopes = append(scopes, input.Scopes[i].Name)
	}
	perm := permission.Permission{
		Name:             permission.NameNonAdmin(input.RoleName, input.ResourceName),
		Description:      fmt.Sprintf("Permission for %s resource with %s role", input.ResourceName, input.RoleName),
		Type:             permission.TypeScope,
		Resources:        []string{input.ResourceID},
		Scopes:           scopes,
		Policies:         []string{p.ID},
		DecisionStrategy: permission.DecisionAffirmative,
	}

	_, err = ks.CreatePermission(ctx, input.TenantDomain, *input.ClientID, perm)
	if err != nil {
		return fmt.Errorf("failed to create permission for %s %s: %w", input.ResourceName, input.RoleName, err)
	}

	return
}
func (ks *KeycloakService) RemoveRolesFromResource(ctx context.Context, input *resourceCmd.ResourceAssignRoleInput) (err error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.fetchClient(ctx, &input.BaseInput)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	permissions, err := ks.client.GetPermissions(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, gocloak.GetPermissionParams{
		Resource: &input.ResourceID,
	})
	for i := range permissions {
		err = ks.client.DeletePermission(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *permissions[i].ID)
		if err != nil {
			return
		}
	}
	return
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
