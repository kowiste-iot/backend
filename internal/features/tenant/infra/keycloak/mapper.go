package keycloak

import (
	"backend/internal/features/tenant/domain/command"
	"backend/shared/util"
	"backend/shared/keycloak"

	"github.com/Nerzal/gocloak/v13"
)



func (ks *BranchKeycloak) convertFromGoCloak(c *gocloak.Client) keycloak.Client {
	if c == nil {
		return keycloak.Client{}
	}
	return keycloak.Client{
		ID:                      c.ID,
		ClientID:                gocloak.PString(c.ClientID),
		Name:                    gocloak.PString(c.Name),
		Description:             gocloak.PString(c.Description),
		RootURL:                 gocloak.PString(c.RootURL),
		AdminURL:                gocloak.PString(c.AdminURL),
		ClientAuthenticatorType: gocloak.PString(c.ClientAuthenticatorType),
		RedirectURIs:            gocloak.PStringSlice(c.RedirectURIs),
		WebOrigins:              gocloak.PStringSlice(c.WebOrigins),
		PublicClient:            gocloak.PBool(c.PublicClient),
		StandardFlowEnabled:     gocloak.PBool(c.StandardFlowEnabled),
		ImplicitFlowEnabled:     gocloak.PBool(c.ImplicitFlowEnabled),
		ServiceAccountEnabled:   gocloak.PBool(c.ServiceAccountsEnabled),
		FullScopeAllowed:        gocloak.PBool(c.FullScopeAllowed),
	}
}

func (ks *BranchKeycloak) convertToGoCloak(isBack bool, branch string) gocloak.Client {
	c := new(keycloak.Client)
	if isBack {
		upperBranch := util.CapitalizeFirst(branch)
		c = &keycloak.Client{
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
		c = &keycloak.Client{
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

	if !c.Authorization {
		attributes["oidc.ciba.grant.enabled"] = "false"
		attributes["post.logout.redirect.uris"] = "http://localhost:5173"
		attributes["display.on.consent.screen"] = "false"
		attributes["oauth2.device.authorization.grant.enabled"] = "false"

		defaultScopes = append(defaultScopes, []string{"acr", "basic"}...)
		optScopes = append(optScopes, "organization")
	}

	data := gocloak.Client{
		ClientID:                     &c.ClientID,
		Enabled:                      gocloak.BoolP(true),
		Description:                  &c.Description,
		ClientAuthenticatorType:      &c.ClientAuthenticatorType,
		RedirectURIs:                 &c.RedirectURIs,
		StandardFlowEnabled:          &c.StandardFlowEnabled,
		DirectAccessGrantsEnabled:    gocloak.BoolP(true),
		PublicClient:                 &c.PublicClient,
		FrontChannelLogout:           gocloak.BoolP(true),
		Protocol:                     gocloak.StringP("openid-connect"),
		Attributes:                   &attributes,
		FullScopeAllowed:             &c.FullScopeAllowed,
		NodeReRegistrationTimeout:    gocloak.Int32P(-1),
		DefaultClientScopes:          &defaultScopes,
		OptionalClientScopes:         &optScopes,
		AuthorizationServicesEnabled: &c.AuthorizationEnabled,
		ServiceAccountsEnabled:       &c.ServiceAccountEnabled,
	}

	if c.ID != nil {
		data.ID = c.ID
	}
	if len(c.WebOrigins) == 0 {
		data.WebOrigins = &[]string{"*"}
	}

	return data
}
