package keycloak

import (
	permissionDomain "backend/internal/features/permission/domain"
	"backend/shared/keycloak"
	roleDomain"backend/internal/features/user/domain"

	baseCmd "backend/shared/base/command"
	"context"
	"fmt"
)



func (ks *BranchKeycloak) createClient(ctx context.Context, isBack bool, input *baseCmd.BaseInput) (*keycloak.Client, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	goClient := ks.convertToGoCloak(isBack, input.BranchName)
	id, err := ks.Client.CreateClient(ctx, token.AccessToken, input.TenantDomain, goClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	if isBack {
		err = ks.updateClientSettings(ctx, ks.Config.Host, token.AccessToken, input.TenantDomain, id, &clientSettings{
			ID:                            id,
			ClientID:                      id,
			Name:                          *goClient.ClientID,
			AllowRemoteResourceManagement: true,
			PolicyEnforcementMode:         roleDomain.Enforcing,
			Resources:                     []string{},
			Policies:                      []string{},
			Scopes:                        []string{},
			DecisionStrategy:              permissionDomain.DecisionAffirmative,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to set authorization client: %w", err)
		}
	}
	createdClient, err := ks.GetClient(ctx, input.TenantDomain, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get created client: %w", err)
	}
	if !isBack {
		ks.createProtocolMapper(ctx, input.TenantDomain, *createdClient.ID)
	}
	return createdClient, nil
}

func (ks *BranchKeycloak) createProtocolMapper(ctx context.Context, tenantDomain, clientID string) (err error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	//add info group and roles to the tokens
	mapper := ProtocolMapper{
		Protocol:       "openid-connect",
		ProtocolMapper: "oidc-group-membership-mapper",
		Name:           "user groups",
		Config: ProtocolMapperConfig{
			ClaimName:               "branch",
			FullPath:                "false",
			IDTokenClaim:            "true",
			AccessTokenClaim:        "true",
			LightweightClaim:        "false",
			UserinfoTokenClaim:      "true",
			IntrospectionTokenClaim: "true",
		},
	}
	err = CreateProtocolMapper(ctx, ks.Config.Host, token.AccessToken, tenantDomain, clientID, mapper)
	if err != nil {
		return fmt.Errorf("failed to set client mapper: %w", err)
	}
	return
}
func (ks *BranchKeycloak) GetClient(ctx context.Context, tenantID, clientID string) (*keycloak.Client, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	client, err := ks.Client.GetClient(ctx, token.AccessToken, tenantID, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	result := ks.convertFromGoCloak(client)
	return &result, nil
}
