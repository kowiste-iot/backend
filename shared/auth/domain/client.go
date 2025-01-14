package auth

import "context"

const (
	BackendClient = "backend-service"
)

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
	BaseURL                   string            `json:"baseUrl"`
	SurrogateAuthRequired     bool              `json:"surrogateAuthRequired"`
	AlwaysDisplayInConsole    bool              `json:"alwaysDisplayInConsole"`
	ClientAuthenticatorType   string            `json:"clientAuthenticatorType"`
	RedirectURIs              []string          `json:"redirectUris"`
	WebOrigins                []string          `json:"webOrigins"`
	NotBefore                 int               `json:"notBefore"`
	BearerOnly                bool              `json:"bearerOnly"`
	ConsentRequired           bool              `json:"consentRequired"`
	StandardFlowEnabled       bool              `json:"standardFlowEnabled"`
	ImplicitFlowEnabled       bool              `json:"implicitFlowEnabled"`
	DirectAccessGrantsEnabled bool              `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled    bool              `json:"serviceAccountsEnabled"`
	PublicClient              bool              `json:"publicClient"`
	Protocol                  string            `json:"protocol"`
	Attributes                map[string]string `json:"attributes"`
	FullScopeAllowed          bool              `json:"fullScopeAllowed"`
	NodeReRegistrationTimeout int               `json:"nodeReRegistrationTimeout"`
	DefaultClientScopes       []string          `json:"defaultClientScopes"`
	OptionalClientScopes      []string          `json:"optionalClientScopes"`
	AuthorizationEnabled      bool              `json:"authorizationEnabled"`
	ServiceAccountEnabled     bool              `json:"serviceAccountEnabled"`
}
