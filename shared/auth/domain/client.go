package auth

import "context"

type ClientProvider interface {
	CreateClient(ctx context.Context, tenantID string, client Client) (*Client, error)
	UpdateClient(ctx context.Context, tenantID string, client Client) error
	DeleteClient(ctx context.Context, tenantID, clientID string) error
	GetClient(ctx context.Context, tenantID, clientID string) (*Client, error)
	GetClientByClientID(ctx context.Context, tenantID, clientID string) (*Client, error)
	ListClients(ctx context.Context, tenantID string) ([]Client, error)
	UpdateClientRoles(ctx context.Context, tenantID, clientID string, roles []string) error
	GetClientRoles(ctx context.Context, tenantID, clientID string) ([]string, error)
}

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
