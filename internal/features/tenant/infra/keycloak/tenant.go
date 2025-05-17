package keycloak

import (
	"backend/internal/features/tenant/domain"
	"backend/pkg/config"
	"context"
	"fmt"
	"time"

	"backend/shared/keycloak"

	"github.com/Nerzal/gocloak/v13"
)

type TenantKeycloak struct {
	*keycloak.Keycloak
	tenantConfig *config.TenantConfiguration
}

func New(cfg *config.TenantConfiguration, core *keycloak.Keycloak) *TenantKeycloak {
	return &TenantKeycloak{
		Keycloak:     core,
		tenantConfig: cfg,
	}
}

func (rk TenantKeycloak) CreateTenant(ctx context.Context, tenant *domain.Tenant) (id string, err error) {

	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	// Create realm
	realm := mapTenantToRealm(tenant)
	id, err = rk.Client.CreateRealm(ctx, token.AccessToken, *realm)
	if err != nil {
		return "", fmt.Errorf("failed to create realm: %w", err)
	}
	//force get token
	_, err = rk.GetNewToken(ctx)
	if err != nil {
		return
	}

	return
}
func (rk TenantKeycloak) DeleteTenant(ctx context.Context, tenantID string) error {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = rk.Client.DeleteRealm(ctx, token.AccessToken, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete realm: %w", err)
	}

	return nil

}
func (rk TenantKeycloak) GetTenant(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Get all realms
	realms, err := rk.Client.GetRealms(ctx, token.AccessToken)
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
func (rk TenantKeycloak) UpdateTenant(ctx context.Context, tenant *domain.Tenant) error {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	realm := mapTenantToRealm(tenant)
	err = rk.Client.UpdateRealm(ctx, token.AccessToken, *realm)
	if err != nil {
		return fmt.Errorf("failed to update realm: %w", err)
	}

	return rk.updateRealmConfig(ctx, tenant)
}

// CreateAdminGroup create a group call admin in the pass tenant
func (rk TenantKeycloak) CreateAdminGroup(ctx context.Context, tenantDomain string) (id string, err error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	group := mapBranchToGroup(domain.AdminBranch, "group for "+domain.AdminBranch)
	id, err = rk.Client.CreateGroup(ctx, token.AccessToken, tenantDomain, *group)
	return
}

func mapTenantToRealm(tenant *domain.Tenant) *gocloak.RealmRepresentation {
	enabled := true

	realm := &gocloak.RealmRepresentation{
		Realm:       gocloak.StringP(tenant.Domain()),
		Enabled:     &enabled,
		DisplayName: gocloak.StringP(tenant.Name()),
	}

	return realm
}
func mapRealmToTenant(realm *gocloak.RealmRepresentation) *domain.Tenant {
	tenant := domain.NewFromRepository(*realm.ID, "", *realm.DisplayName, *realm.Realm, "", time.Now(), nil)
	return tenant
}
func (ks *TenantKeycloak) updateRealmConfig(ctx context.Context, tenant *domain.Tenant) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	clientID := fmt.Sprintf("%s-client", tenant.Domain())
	clients, err := ks.Client.GetClients(ctx, token.AccessToken, tenant.Domain(), gocloak.GetClientsParams{
		ClientID: &clientID,
	})
	if err != nil {
		return fmt.Errorf("failed to get clients: %w", err)
	}

	if len(clients) > 0 {
		client := clients[0]

		client.StandardFlowEnabled = gocloak.BoolP(true)

		err = ks.Client.UpdateClient(ctx, token.AccessToken, tenant.Domain(), *client)
		if err != nil {
			return fmt.Errorf("failed to update client: %w", err)
		}
	}

	return nil
}
