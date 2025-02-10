package keycloak

import (
	"backend/shared/auth/domain/scope"
	baseCmd "backend/shared/base/command"
	"context"
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
	}, nil
}


func (ks *KeycloakService) ListScopes(ctx context.Context, input *baseCmd.BaseInput) ([]scope.Scope, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.fetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	kcScopes, err := ks.client.GetScopes(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
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
		}
	}

	return scopes, nil
}


