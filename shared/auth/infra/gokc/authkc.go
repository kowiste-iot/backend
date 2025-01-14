package keycloak

import (
	"context"
	"ddd/shared/http/httputil"
	"errors"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
)

type KeycloakConfig struct {
	Host         string
	Realm        string
	ClientID     string
	ClientSecret string
	WebClientID  string
}

type KeycloakService struct {
	client      *gocloak.GoCloak
	config      KeycloakConfig
	realmToken  *gocloak.JWT
	mutex       sync.RWMutex
	lastRefresh time.Time
}

func NewKeycloakService(config KeycloakConfig) (*KeycloakService, error) {

	ks := &KeycloakService{
		client: gocloak.NewClient(config.Host),
		config: config,
	}

	// Initialize token and resources
	if err := ks.initialize(context.Background()); err != nil {
		return nil, err
	}
	return ks, nil
}
func (ks *KeycloakService) initialize(ctx context.Context) error {
	// Get initial token
	token, err := ks.getNewToken(ctx)
	if err != nil {
		return err
	}

	ks.mutex.Lock()
	ks.realmToken = token
	ks.mutex.Unlock()

	// Load resources
	// if err := ks.refreshResources(ctx); err != nil {
	// 	return err
	// }

	return nil
}
func (ks *KeycloakService) getNewToken(ctx context.Context) (token *gocloak.JWT, err error) {
	token, err = ks.client.GetToken(ctx, ks.config.Realm, gocloak.TokenOptions{
		ClientID:     &ks.config.ClientID,
		ClientSecret: &ks.config.ClientSecret,
		GrantType:    gocloak.StringP("client_credentials"),
	})
	if err != nil {
		return
	}
	ks.realmToken = token
	ks.lastRefresh = time.Now()
	return
}

func (ks *KeycloakService) isTokenExpired() bool {
	if ks.realmToken == nil {
		return true
	}
	elapsed := time.Since(ks.lastRefresh).Seconds()
	remaining := float64(ks.realmToken.ExpiresIn) - elapsed
	return remaining <= 10
}
func (ks *KeycloakService) GetValidToken(ctx context.Context) (*gocloak.JWT, error) {
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
	newToken, err := ks.getNewToken(ctx)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// DecodeAccessToken decodes and validates the access token
func (k *KeycloakService) ValidateToken(ctx context.Context, accessToken string) (*jwt.Token, error) {
	tenant := k.getTenantOrDefault(ctx)
	decodedToken, _, err := k.client.DecodeAccessToken(
		ctx,
		accessToken,
		tenant,
	)

	if err != nil {
		return nil, err
	}

	if !decodedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return &jwt.Token{
		Raw:       decodedToken.Raw,
		Method:    decodedToken.Method,
		Header:    decodedToken.Header,
		Claims:    decodedToken.Claims,
		Signature: decodedToken.Signature,
		Valid:     decodedToken.Valid,
	}, nil
}

func (k *KeycloakService) getTenantOrDefault(ctx context.Context) string {
	domain := k.config.Realm
	tenant, ok := httputil.GetTenant(ctx)
	if ok {
		domain = tenant.AuhtID()
	}
	return domain
}
