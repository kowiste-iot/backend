package keycloak

import (

	baseCmd "backend/shared/base/command"
	"context"
	"fmt"
	tenantCmd "backend/internal/features/tenant/domain/command"

	"github.com/Nerzal/gocloak/v13"
)

type Client struct {
	ID                        *string           `json:"id"`
	ClientID                  string            `json:"clientId"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	RootURL                   string            `json:"rootUrl"`
	AdminURL                  string            `json:"adminUrl"`
	ClientAuthenticatorType   string            `json:"clientAuthenticatorType"`
	RedirectURIs              []string          `json:"redirectUris"`
	WebOrigins                []string          `json:"webOrigins"`
	StandardFlowEnabled       bool              `json:"standardFlowEnabled"`
	ImplicitFlowEnabled       bool              `json:"implicitFlowEnabled"`
	PublicClient              bool              `json:"publicClient"`
	FullScopeAllowed          bool              `json:"fullScopeAllowed"`
	AuthorizationEnabled      bool              `json:"authorizationEnabled"`
	ServiceAccountEnabled     bool              `json:"serviceAccountEnabled"`
	Authorization             bool              `json:"authorization"`
}

func (ks *Keycloak) GetClientByClientID(ctx context.Context, tenantID, clientID string) (*Client, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	clients, err := ks.Client.GetClients(ctx, token.AccessToken, tenantID, gocloak.GetClientsParams{
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
func (ks *Keycloak) convertToGoCloak(client Client) gocloak.Client {
	attributes := map[string]string{
		"realm_client":                             "false",
		"backchannel.logout.session.required":      "true",
		"backchannel.logout.revoke.offline.tokens": "false",
	}

	defaultScopes := []string{
		"web-origins",
		"roles",
		"profile",
		"email",
	}
	optScopes := []string{
		"address",
		"phone",
		"offline_access",
		"microprofile-jwt",
	}

	if !client.Authorization {
		attributes["oidc.ciba.grant.enabled"] = "false"
		attributes["post.logout.redirect.uris"] = "http://localhost:5173"
		attributes["display.on.consent.screen"] = "false"
		attributes["oauth2.device.authorization.grant.enabled"] = "false"

		defaultScopes = append(defaultScopes, []string{"acr", "basic"}...)
		optScopes = append(optScopes, "organization")
	}

	data := gocloak.Client{
		ClientID:                     &client.ClientID,
		Enabled:                      gocloak.BoolP(true),
		Description:                  &client.Description,
		ClientAuthenticatorType:      &client.ClientAuthenticatorType,
		RedirectURIs:                 &client.RedirectURIs,
		StandardFlowEnabled:          &client.StandardFlowEnabled,
		DirectAccessGrantsEnabled:    gocloak.BoolP(true),
		PublicClient:                 &client.PublicClient,
		FrontChannelLogout:           gocloak.BoolP(true),
		Protocol:                     gocloak.StringP("openid-connect"),
		Attributes:                   &attributes,
		FullScopeAllowed:             &client.FullScopeAllowed,
		NodeReRegistrationTimeout:    gocloak.Int32P(-1),
		DefaultClientScopes:          &defaultScopes,
		OptionalClientScopes:         &optScopes,
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

func (ks *Keycloak) convertFromGoCloak(client *gocloak.Client) Client {
	if client == nil {
		return Client{}
	}
	return Client{
		ID:                      client.ID,
		ClientID:                gocloak.PString(client.ClientID),
		Name:                    gocloak.PString(client.Name),
		Description:             gocloak.PString(client.Description),
		RootURL:                 gocloak.PString(client.RootURL),
		AdminURL:                gocloak.PString(client.AdminURL),
		ClientAuthenticatorType: gocloak.PString(client.ClientAuthenticatorType),
		RedirectURIs:            gocloak.PStringSlice(client.RedirectURIs),
		WebOrigins:              gocloak.PStringSlice(client.WebOrigins),
		PublicClient:            gocloak.PBool(client.PublicClient),
		StandardFlowEnabled:     gocloak.PBool(client.StandardFlowEnabled),
		ImplicitFlowEnabled:     gocloak.PBool(client.ImplicitFlowEnabled),
		ServiceAccountEnabled:   gocloak.PBool(client.ServiceAccountsEnabled),
		FullScopeAllowed:        gocloak.PBool(client.FullScopeAllowed),
	}
}
func (k *Keycloak) FetchClient(ctx context.Context, input *baseCmd.BaseInput) (err error) {
	if input.ClientID == nil {
		client, err := k.GetClientByClientID(ctx, input.TenantDomain, tenantCmd.ClientName(input.BranchName))
		if err != nil {
			return fmt.Errorf("error getting client: %w", err)
		}
		input.ClientID = client.ID
	}
	return
}
