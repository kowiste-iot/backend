package keycloak

import (
	"backend/internal/features/scope/domain"
	baseCmd "backend/shared/base/command"
	"context"
	"fmt"

	"backend/shared/keycloak"

	"github.com/Nerzal/gocloak/v13"
)

type ScopeKeycloak struct {
	*keycloak.Keycloak
}

func New(core *keycloak.Keycloak) *ScopeKeycloak {
	return &ScopeKeycloak{
		Keycloak: core,
	}
}

func (ks *ScopeKeycloak) CreateScope(ctx context.Context, input *baseCmd.BaseInput, sc domain.Scope) (*domain.Scope, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	kcScope := gocloak.ScopeRepresentation{
		Name:        &sc.Name,
		DisplayName: &sc.DisplayName,
	}
	err = ks.FetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	createdScope, err := ks.Client.CreateScope(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		kcScope,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scope: %w", err)
	}

	return &domain.Scope{
		ID:          *createdScope.ID,
		Name:        *createdScope.Name,
		DisplayName: *createdScope.DisplayName,
	}, nil
}

func (rk ScopeKeycloak) ListScopes(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Scope, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	kcScopes, err := rk.Client.GetScopes(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		gocloak.GetScopeParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list scopes: %w", err)
	}

	scopes := make([]domain.Scope, len(kcScopes))
	for i, ks := range kcScopes {
		scopes[i] = domain.Scope{
			ID:          *ks.ID,
			Name:        *ks.Name,
			DisplayName: *ks.DisplayName,
		}
	}

	return scopes, nil
}
