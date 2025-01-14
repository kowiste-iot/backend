package keycloak

import (
	"context"
	"ddd/shared/auth/domain/scope"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

func (ks *KeycloakService) CreateScope(ctx context.Context, tenantID, clientID string, sc scope.Scope) (*scope.Scope, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	kcScope := gocloak.ScopeRepresentation{
		Name:        &sc.Name,
		DisplayName: &sc.DisplayName,
		IconURI:     &sc.IconURI,
	}

	createdScope, err := ks.client.CreateScope(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		kcScope,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scope: %w", err)
	}

	return &scope.Scope{
		ID:          *createdScope.ID,
		Name:        *createdScope.Name,
		DisplayName: *createdScope.DisplayName,
		IconURI:     *createdScope.IconURI,
	}, nil
}
func (ks *KeycloakService) GetScope(ctx context.Context, tenantID, clientID, scopeID string) (*scope.Scope, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	kcScope, err := ks.client.GetScope(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		scopeID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get scope: %w", err)
	}

	return &scope.Scope{
		ID:          *kcScope.ID,
		Name:        *kcScope.Name,
		DisplayName: *kcScope.DisplayName,
		IconURI:     *kcScope.IconURI,
	}, nil
}

func (ks *KeycloakService) ListScopes(ctx context.Context, tenantID, clientID string) ([]scope.Scope, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	kcScopes, err := ks.client.GetScopes(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		gocloak.GetScopeParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list scopes: %w", err)
	}

	scopes := make([]scope.Scope, len(kcScopes))
	for i, ks := range kcScopes {
		scopes[i] = scope.Scope{
			ID:          *ks.ID,
			Name:        *ks.Name,
			DisplayName: *ks.DisplayName,
			IconURI:     *ks.IconURI,
		}
	}

	return scopes, nil
}
func (ks *KeycloakService) UpdateScope(ctx context.Context, tenantID, clientID string, scope scope.Scope) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	kcScope := gocloak.ScopeRepresentation{
		ID:          &scope.ID,
		Name:        &scope.Name,
		DisplayName: &scope.DisplayName,
		IconURI:     &scope.IconURI,
	}

	err = ks.client.UpdateScope(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		kcScope,
	)
	if err != nil {
		return fmt.Errorf("failed to update scope: %w", err)
	}

	return nil
}

func (ks *KeycloakService) DeleteScope(ctx context.Context, tenantID, clientID, scopeID string) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.DeleteScope(
		ctx,
		token.AccessToken,
		tenantID,
		clientID,
		scopeID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete scope: %w", err)
	}

	return nil
}
