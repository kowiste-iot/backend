package keycloak

import (
	"backend/shared/http/httputil"
	"context"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"
)

type Keycloak struct {
	Client      *gocloak.GoCloak
	Config      *KeycloakConfig
	realmToken  *gocloak.JWT
	mutex       sync.RWMutex
	lastRefresh time.Time
}
type KeycloakConfig struct {
	Host         string
	Realm        string
	ClientID     string
	ClientSecret string
	WebClientID  string
}

func New(config *KeycloakConfig) (*Keycloak, error) {
	kc := new(Keycloak)
	kc.Client = gocloak.NewClient(config.Host)
	kc.Config = config
	if err := kc.initialize(context.Background()); err != nil {
		return nil, err
	}
	return kc, nil
}

func (ks *Keycloak) initialize(ctx context.Context) error {
	// Get initial token
	token, err := ks.GetNewToken(ctx)
	if err != nil {
		return err
	}

	ks.mutex.Lock()
	ks.realmToken = token
	ks.mutex.Unlock()

	return nil
}
func (ks *Keycloak) GetNewToken(ctx context.Context) (token *gocloak.JWT, err error) {
	token, err = ks.Client.GetToken(ctx, ks.Config.Realm, gocloak.TokenOptions{
		ClientID:     &ks.Config.ClientID,
		ClientSecret: &ks.Config.ClientSecret,
		GrantType:    gocloak.StringP("client_credentials"),
	})
	if err != nil {
		return
	}
	ks.realmToken = token
	ks.lastRefresh = time.Now()
	return
}
func (ks *Keycloak) isTokenExpired() bool {
	if ks.realmToken == nil {
		return true
	}
	elapsed := time.Since(ks.lastRefresh).Seconds()
	remaining := float64(ks.realmToken.ExpiresIn) - elapsed
	return remaining <= 10
}
func (ks *Keycloak) GetValidToken(ctx context.Context) (*gocloak.JWT, error) {
	ks.mutex.RLock()
	if ks.realmToken != nil && !ks.isTokenExpired() {
		token := ks.realmToken
		ks.mutex.RUnlock()
		return token, nil
	}
	ks.mutex.RUnlock()

	ks.mutex.Lock()
	defer ks.mutex.Unlock()

	// Get new token
	newToken, err := ks.GetNewToken(ctx)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

func (k *Keycloak) getTenantOrDefault(ctx context.Context) string {
	domain := k.Config.Realm
	tenant, ok := httputil.GetTenant(ctx)
	if ok {
		domain = tenant.Domain()
	}
	return domain
}
