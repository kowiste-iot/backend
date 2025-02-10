package keycloak

import (
	"backend/internal/features/tenant/domain/command"
	auth "backend/shared/auth/domain"
	"backend/shared/util"

	"github.com/Nerzal/gocloak/v13"
)

func (ks *BranchKeycloak) convertFromGoCloak(client *gocloak.Client) auth.Client {
	if client == nil {
		return auth.Client{}
	}
	return auth.Client{
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

func (ks *BranchKeycloak) convertToGoCloak2(isBack bool, branch string) gocloak.Client {
	client := new(auth.Client)
	if isBack {
		upperBranch := util.CapitalizeFirst(branch)
		client = &auth.Client{
			ClientID:                command.ClientName(branch),
			Name:                    upperBranch + ks.tenantConfig.BackendClient.Name,
			Description:             upperBranch + ks.tenantConfig.BackendClient.Description,
			RootURL:                 *ks.tenantConfig.BackendClient.RootURL,
			AdminURL:                *ks.tenantConfig.BackendClient.AdminURL,
			ClientAuthenticatorType: "client-secret",
			RedirectURIs:            ks.tenantConfig.BackendClient.RedirectURIs,
			WebOrigins:              ks.tenantConfig.BackendClient.Origins,
			ServiceAccountEnabled:   true,
			AuthorizationEnabled:    true,
		}
	} else {
		client = &auth.Client{
			ClientID:            ks.tenantConfig.WebClient.ClientID,
			Name:                ks.tenantConfig.WebClient.Name,
			RedirectURIs:        ks.tenantConfig.WebClient.RedirectURIs,
			WebOrigins:          ks.tenantConfig.WebClient.Origins,
			StandardFlowEnabled: true,
			PublicClient:        true,
		}
	}
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
