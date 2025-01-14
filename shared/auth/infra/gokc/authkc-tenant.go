package keycloak

import (
	"context"
	auth "ddd/shared/auth/domain"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

// TenantProvider implementation using the existing KeycloakService
func (ks *KeycloakService) CreateTenant(ctx context.Context, tenant *auth.Tenant) (id string, err error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	// Create realm
	realm := mapTenantToRealm(tenant)
	id, err = ks.client.CreateRealm(ctx, token.AccessToken, *realm)
	if err != nil {
		return "", fmt.Errorf("failed to create realm: %w", err)
	}
	//force get token
	_, err = ks.getNewToken(ctx)
	if err != nil {
		return
	}
	return
}

func (ks *KeycloakService) UpdateTenant(ctx context.Context, tenant *auth.Tenant) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	realm := mapTenantToRealm(tenant)
	err = ks.client.UpdateRealm(ctx, token.AccessToken, *realm)
	if err != nil {
		return fmt.Errorf("failed to update realm: %w", err)
	}

	return ks.updateRealmConfig(ctx, tenant)
}

func (ks *KeycloakService) DeleteTenant(ctx context.Context, tenantDomain string) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.DeleteRealm(ctx, token.AccessToken, tenantDomain)
	if err != nil {
		return fmt.Errorf("failed to delete realm: %w", err)
	}

	return nil
}

func (ks *KeycloakService) GetTenant(ctx context.Context, tenantID string) (*auth.Tenant, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Get all realms
	realms, err := ks.client.GetRealms(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get realms: %w", err)
	}

	// Find realm by tenant ID in attributes
	for _, realm := range realms {
		if *realm.Realm == tenantID {
			return mapRealmToTenant(realm), nil
		}
	}

	return nil, fmt.Errorf("tenant not found: %s", tenantID)
}

// Helper functions

func mapTenantToRealm(tenant *auth.Tenant) *gocloak.RealmRepresentation {
	enabled := true
	realm := &gocloak.RealmRepresentation{
		Realm:       &tenant.Domain,
		Enabled:     &enabled,
		DisplayName: &tenant.Name,
		LoginTheme:  &tenant.Theme.Login,
	}

	return realm
}

func mapRealmToTenant(realm *gocloak.RealmRepresentation) *auth.Tenant {
	tenant := &auth.Tenant{
		ID:     *realm.ID,
		Domain: *realm.Realm,
		Name:   *realm.DisplayName,
		Theme:  &auth.TenantTheme{
			Login: *realm.LoginTheme,
		},
	}

	return tenant
}

func (ks *KeycloakService) setupRealmConfig(ctx context.Context, tenant auth.Tenant) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Create default client
	clientID := fmt.Sprintf("%s-client", tenant.Domain)
	client := &gocloak.Client{
		ClientID:                  &clientID,
		Enabled:                   gocloak.BoolP(true),
		StandardFlowEnabled:       gocloak.BoolP(true),
		DirectAccessGrantsEnabled: gocloak.BoolP(true),
		Protocol:                  gocloak.StringP("openid-connect"),
	}

	_, err = ks.client.CreateClient(ctx, token.AccessToken, tenant.Domain, *client)
	if err != nil {
		return fmt.Errorf("failed to create default client: %w", err)
	}

	// Create default roles
	defaultRoles := []string{"user", "admin"}
	for _, roleName := range defaultRoles {
		role := gocloak.Role{
			Name: &roleName,
		}
		_, err := ks.client.CreateRealmRole(ctx, token.AccessToken, tenant.Domain, role)
		if err != nil {
			return fmt.Errorf("failed to create role %s: %w", roleName, err)
		}
	}

	return nil
}

func (ks *KeycloakService) updateRealmConfig(ctx context.Context, tenant *auth.Tenant) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	clientID := fmt.Sprintf("%s-client", tenant.Domain)
	clients, err := ks.client.GetClients(ctx, token.AccessToken, tenant.Domain, gocloak.GetClientsParams{
		ClientID: &clientID,
	})
	if err != nil {
		return fmt.Errorf("failed to get clients: %w", err)
	}

	if len(clients) > 0 {
		client := clients[0]
		// Update client settings based on tenant configuration
		if tenant.Settings.Features["oauth_enabled"] {
			client.StandardFlowEnabled = gocloak.BoolP(true)
		} else {
			client.StandardFlowEnabled = gocloak.BoolP(false)
		}

		err = ks.client.UpdateClient(ctx, token.AccessToken, tenant.Domain, *client)
		if err != nil {
			return fmt.Errorf("failed to update client: %w", err)
		}
	}

	return nil
}
