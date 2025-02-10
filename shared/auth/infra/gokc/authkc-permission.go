package keycloak

import (
	"backend/shared/auth/domain/permission"
	"backend/shared/auth/infra/restkc"
	baseCmd "backend/shared/base/command"
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

// ValidatePermissionService checks if the service/token has permission to access a resource with specific scope
func (k *KeycloakService) ValidatePermissionService(ctx context.Context, token, clientID, resource, scope string) (bool, error) {
	tenant := k.getTenantOrDefault(ctx)
	permissions, err := k.client.GetRequestingPartyPermissions(
		ctx,
		token,
		tenant,
		gocloak.RequestingPartyTokenOptions{
			GrantType:   gocloak.StringP("urn:ietf:params:oauth:grant-type:uma-ticket"),
			Audience:    &clientID,
			Permissions: &[]string{resource},
		},
	)

	if err != nil {
		return false, err
	}

	// Handle case where permissions is nil
	if permissions == nil {
		return false, nil
	}

	// Check if the user has the required scope
	for _, permission := range *permissions {
		if permission.ResourceID != nil && *permission.ResourceID == resource {
			if permission.Scopes != nil {
				for _, s := range *permission.Scopes {
					if s == scope {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

// ValidatePermissionUser checks if the user has permission to access a resource with specific scope
func (k *KeycloakService) ValidatePermissionUser(ctx context.Context, token, clientID, resource, action string) (bool, error) {
	tenant := k.getTenantOrDefault(ctx)
	permissions := []string{resource}

	result, err := k.client.GetRequestingPartyPermissions(ctx,
		token,
		tenant,
		gocloak.RequestingPartyTokenOptions{
			Permissions: &permissions,
			Audience:    &clientID,
		})
	if err != nil {
		return false, err
	}
	hasAccess := false
	if permissions != nil {
		for _, permission := range *result {
			if permission.ResourceName != nil && *permission.ResourceName == resource {
				if permission.Scopes == nil {
					// If scopes is nil means full access
					hasAccess = true
					break
				}
				for _, scope := range *permission.Scopes {
					if scope == action {
						hasAccess = true
						break
					}
				}
			}
		}
	}
	return hasAccess, nil

}

func (ks *KeycloakService) CreatePermission(ctx context.Context, tenantDomain, clientID string, p permission.Permission) (*permission.Permission, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	internalScopes := []string{}
	scopes, err := ks.ListScopes(ctx, &baseCmd.BaseInput{TenantDomain: tenantDomain, ClientID: &clientID})
	if err != nil {
		return nil, err
	}

	// Use the scopes list for all iterations
	for _, scopeName := range p.Scopes {
		for _, scope := range scopes {
			if scope.Name == scopeName {
				internalScopes = append(internalScopes, scope.ID)
				break
			}
		}
	}
	p.Scopes = internalScopes

	c, err := ks.client.GetClient(ctx, token.AccessToken, tenantDomain, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission error client: %w", err)
	}

	created, err := restkc.CreatePermission(ctx, ks.config.Host, token.AccessToken, tenantDomain, *c.ID, p)

	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	// Convert back to our Permission type
	result := &permission.Permission{
		ID:               created.ID,
		Name:             created.Name,
		DecisionStrategy: created.DecisionStrategy,
		Description:      created.Description,
		Resources:        created.Resources,
		Scopes:           created.Scopes,
		Policies:         created.Policies,
	}

	return result, nil
}

func (ks *KeycloakService) GetPermission(ctx context.Context, tenantID, clientID, permissionID string) (*permission.Permission, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	kcPermission, err := ks.client.GetPermission(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		permissionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &permission.Permission{
		ID:        *kcPermission.ID,
		Name:      *kcPermission.Name,
		Type:      *kcPermission.Type,
		Resources: *kcPermission.Resources,
		Scopes:    *kcPermission.Scopes,
		Policies:  *kcPermission.Policies,
	}, nil
}

func (ks *KeycloakService) ListPermissions(ctx context.Context, input *baseCmd.BaseInput) ([]permission.Permission, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.fetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	kcPermissions, err := ks.client.GetPermissions(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		gocloak.GetPermissionParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	permissions := make([]permission.Permission, len(kcPermissions))
	for i, kp := range kcPermissions {
		permissions[i] = permission.Permission{
			ID:    *kp.ID,
			Name:  *kp.Name,
			Type:  *kp.Type,
			Logic: string(*kp.Logic),
		}
		if kp.Resources != nil {
			permissions[i].Resources = *kp.Resources
		}
		if kp.Scopes != nil {
			permissions[i].Scopes = *kp.Scopes
		}
		if kp.Resources != nil {
			permissions[i].Policies = *kp.Policies
		}
		resourPerm, err := ks.client.GetPermissionResources(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *kp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch permission resource %w", err)
		}
		if len(resourPerm) > 0 {
			for j := range resourPerm {
				permissions[i].Resources = append(permissions[i].Resources, *resourPerm[j].ResourceID)
			}
		}
		scopesPerm, err := ks.client.GetPermissionScopes(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *kp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch permission scopes %w", err)
		}
		if len(scopesPerm) > 0 {
			for j := range scopesPerm {
				permissions[i].Scopes = append(permissions[i].Scopes, *scopesPerm[j].ScopeName)
			}
		}
		policyPerm, err := ks.client.GetAuthorizationPolicyAssociatedPolicies(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *kp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch permission scopes %w", err)
		}
		if len(policyPerm) > 0 {
			for j := range policyPerm {
				permissions[i].Policies = append(permissions[i].Policies, *policyPerm[j].Name)
			}
		}
	}
	return permissions, nil
}
func (ks *KeycloakService) UpdatePermission(ctx context.Context, tenantID, clientID string, permission permission.Permission) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	kcPermission := gocloak.PermissionRepresentation{
		ID:        &permission.ID,
		Name:      &permission.Name,
		Type:      &permission.Type,
		Resources: &permission.Resources,
		Scopes:    &permission.Scopes,
		Policies:  &permission.Policies,
	}

	err = ks.client.UpdatePermission(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		kcPermission,
	)
	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	return nil
}
func (ks *KeycloakService) DeletePermission(ctx context.Context, tenantID, clientID, permissionID string) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.DeletePermission(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		permissionID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}
