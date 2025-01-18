package keycloak

import (
	"context"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/policy"
	"ddd/shared/auth/infra/restkc"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

func (ks *KeycloakService) CreateClient(ctx context.Context, tenantDomain string, client auth.Client) (*auth.Client, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	goClient := ks.convertToGoCloak(client)
	id, err := ks.client.CreateClient(ctx, token.AccessToken, tenantDomain, goClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	if client.Authorization {
		err = restkc.UpdateClientSettings(ctx, ks.config.Host, token.AccessToken, tenantDomain, id, &restkc.ClientSettings{
			ID:                            id,
			ClientID:                      id,
			Name:                          client.ClientID,
			AllowRemoteResourceManagement: true,
			PolicyEnforcementMode:         policy.Enforcing,
			Resources:                     []string{},
			Policies:                      []string{},
			Scopes:                        []string{},
			DecisionStrategy:              permission.DecisionAffirmative,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to set authorization client: %w", err)
		}
	}
	createdClient, err := ks.GetClient(ctx, tenantDomain, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get created client: %w", err)
	}
	if !client.Authorization {
		//add info group and roles to the tokens
		mapper := restkc.ProtocolMapper{
			Protocol:       "openid-connect",
			ProtocolMapper: "oidc-group-membership-mapper",
			Name:           "user groups",
			Config: restkc.ProtocolMapperConfig{
				ClaimName:               "branch",
				FullPath:                "false",
				IDTokenClaim:            "true",
				AccessTokenClaim:        "true",
				LightweightClaim:        "false",
				UserinfoTokenClaim:      "true",
				IntrospectionTokenClaim: "true",
			},
		}
		err = restkc.CreateProtocolMapper(ctx, ks.config.Host, token.AccessToken, tenantDomain, *createdClient.ID, mapper)
		if err != nil {
			return nil, fmt.Errorf("failed to set client mapper: %w", err)
		}
	}
	return createdClient, nil
}

func (ks *KeycloakService) UpdateClient(ctx context.Context, tenantID string, client auth.Client) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	goClient := ks.convertToGoCloak(client)
	err = ks.client.UpdateClient(ctx, token.AccessToken, tenantID, goClient)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	return nil
}

func (ks *KeycloakService) DeleteClient(ctx context.Context, tenantID, clientID string) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.client.DeleteClient(ctx, token.AccessToken, tenantID, clientID)
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}

func (ks *KeycloakService) GetClient(ctx context.Context, tenantID, clientID string) (*auth.Client, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	client, err := ks.client.GetClient(ctx, token.AccessToken, tenantID, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	result := ks.convertFromGoCloak(client)
	return &result, nil
}

func (ks *KeycloakService) GetClientByClientID(ctx context.Context, tenantID, clientID string) (*auth.Client, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	clients, err := ks.client.GetClients(ctx, token.AccessToken, tenantID, gocloak.GetClientsParams{
		ClientID: &clientID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get client by clientID: %w", err)
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("client not found")
	}

	result := ks.convertFromGoCloak(clients[0])
	return &result, nil
}

func (ks *KeycloakService) ListClients(ctx context.Context, tenantID string) ([]auth.Client, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	clients, err := ks.client.GetClients(ctx, token.AccessToken, tenantID, gocloak.GetClientsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}

	var result []auth.Client
	for _, client := range clients {
		result = append(result, ks.convertFromGoCloak(client))
	}

	return result, nil
}

func (ks *KeycloakService) UpdateClientRoles(ctx context.Context, tenantID, clientID string, roles []string) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	// First, get existing roles
	existingRoles, err := ks.client.GetClientRoles(ctx, token.AccessToken, tenantID, clientID, gocloak.GetRoleParams{})
	if err != nil {
		return fmt.Errorf("failed to get client roles: %w", err)
	}

	// Delete existing roles
	for _, role := range existingRoles {
		err = ks.client.DeleteClientRole(ctx, token.AccessToken, tenantID, clientID, *role.Name)
		if err != nil {
			return fmt.Errorf("failed to delete client role: %w", err)
		}
	}

	// Create new roles
	for _, role := range roles {
		_, err = ks.client.CreateClientRole(ctx, token.AccessToken, tenantID, clientID, gocloak.Role{
			Name: &role,
		})
		if err != nil {
			return fmt.Errorf("failed to create client role: %w", err)
		}
	}

	return nil
}

func (ks *KeycloakService) GetClientRoles(ctx context.Context, tenantID, clientID string) ([]string, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	roles, err := ks.client.GetClientRoles(ctx, token.AccessToken, tenantID, clientID, gocloak.GetRoleParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get client roles: %w", err)
	}

	var result []string
	for _, role := range roles {
		if role.Name != nil {
			result = append(result, *role.Name)
		}
	}

	return result, nil
}

func getStringMap(m *map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return *m
}
func (ks *KeycloakService) convertToGoCloak(client auth.Client) gocloak.Client {
	data := gocloak.Client{
		ClientID:                     &client.ClientID,
		Enabled:                      gocloak.BoolP(true),
		Description:                  &client.Description,
		ClientAuthenticatorType:      &client.ClientAuthenticatorType,
		RedirectURIs:                 &client.RedirectURIs,
		StandardFlowEnabled:          &client.StandardFlowEnabled,
		DirectAccessGrantsEnabled:    &client.DirectAccessGrantsEnabled,
		PublicClient:                 &client.PublicClient,
		FrontChannelLogout:           gocloak.BoolP(true),
		Protocol:                     &client.Protocol,
		Attributes:                   &client.Attributes,
		FullScopeAllowed:             &client.FullScopeAllowed,
		NodeReRegistrationTimeout:    gocloak.Int32P(-1),
		DefaultClientScopes:          &client.DefaultClientScopes,
		OptionalClientScopes:         &client.OptionalClientScopes,
		AuthorizationServicesEnabled: &client.AuthorizationEnabled,
		ServiceAccountsEnabled:       &client.ServiceAccountEnabled,
	}

	if client.ID != nil {
		data.ID = client.ID
	}
	if len(client.WebOrigins) == 0 {
		data.WebOrigins = &[]string{"*"}
	}

	return data
}

func (ks *KeycloakService) convertFromGoCloak(client *gocloak.Client) auth.Client {
	if client == nil {
		return auth.Client{}
	}
	return auth.Client{
		ID:                        client.ID,
		ClientID:                  gocloak.PString(client.ClientID),
		Name:                      gocloak.PString(client.Name),
		Description:               gocloak.PString(client.Description),
		RootURL:                   gocloak.PString(client.RootURL),
		AdminURL:                  gocloak.PString(client.AdminURL),
		BaseURL:                   gocloak.PString(client.BaseURL),
		ClientAuthenticatorType:   gocloak.PString(client.ClientAuthenticatorType),
		RedirectURIs:              gocloak.PStringSlice(client.RedirectURIs),
		WebOrigins:                gocloak.PStringSlice(client.WebOrigins),
		Protocol:                  gocloak.PString(client.Protocol),
		PublicClient:              gocloak.PBool(client.PublicClient),
		Attributes:                getStringMap(client.Attributes),
		DefaultClientScopes:       gocloak.PStringSlice(client.DefaultClientScopes),
		OptionalClientScopes:      gocloak.PStringSlice(client.OptionalClientScopes),
		DirectAccessGrantsEnabled: gocloak.PBool(client.DirectAccessGrantsEnabled),
		StandardFlowEnabled:       gocloak.PBool(client.StandardFlowEnabled),
		ImplicitFlowEnabled:       gocloak.PBool(client.ImplicitFlowEnabled),
		ServiceAccountsEnabled:    gocloak.PBool(client.ServiceAccountsEnabled),
		FullScopeAllowed:          gocloak.PBool(client.FullScopeAllowed),
	}
}

func (k *KeycloakService) getClientToken(ctx context.Context, tenant, clientID string) (token string, err error) {
	cl, err := k.GetClientByClientID(ctx, tenant, clientID)
	if err != nil {
		return "", fmt.Errorf("failed to get client: %w", err)
	}
	secr, err := k.client.GetClientSecret(ctx, k.realmToken.AccessToken, tenant, *cl.ID)
	if err != nil {
		return "s", fmt.Errorf("failed to get user info: %w", err)
	}
	t, err := k.client.LoginClient(ctx, clientID, *secr.Value, tenant)

	return t.AccessToken, nil
}
