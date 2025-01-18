package resource

import (
	"context"
	"ddd/shared/auth/domain/scope"
)

const (
	//Endpoints resources
	Asset  string = "asset-resource"
	Tenant string = "tenant-resource"
	Branch string = "branch-resource"
	User   string = "user-resource"
	Role   string = "role-resource"
	Admin  string = "admin-resource"
	//Endpoin type
	TypeBase = "base-type"
)

type ResourceProvider interface {
	CreateResource(ctx context.Context, tenantID, clientID string, resource Resource) (*Resource, error)
	UpdateResource(ctx context.Context, tenantID, clientID string, resource Resource) error
	DeleteResource(ctx context.Context, tenantID, clientID, resourceID string) error
	GetResource(ctx context.Context, tenantID, clientID, resourceID string) (*Resource, error)
	ListResources(ctx context.Context, tenantID, clientID string) ([]Resource, error)
}

type Resource struct {
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name"`
	Type        string              `json:"type,omitempty"`
	URIs        []string            `json:"uris,omitempty"`
	Scopes      []string            `json:"scopes,omitempty"`
	Attributes  map[string][]string `json:"attributes,omitempty"`
	DisplayName string              `json:"displayName,omitempty"`
	IconURI     string              `json:"icon_uri,omitempty"`
}

func EndpointsResources() (resources []Resource) {
	return []Resource{
		{
			Name:   Asset,
			Type:   TypeBase,
			Scopes: []string{scope.View, scope.Create, scope.Update, scope.Delete},
		},
		{
			Name:   User,
			Type:   TypeBase,
			Scopes: []string{scope.View, scope.Create, scope.Update, scope.Delete},
		},
		{
			Name:   Admin,
			Type:   TypeBase,
			Scopes: []string{scope.View, scope.Create, scope.Update, scope.Delete},
		},
	}
}
