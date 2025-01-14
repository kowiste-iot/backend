package resource

import (
	"context"
)

const (
	Asset  = "asset-management"
	Tenant = "tenant-management"
	Branch = "branch-management"
	User   = "user-management"
	Role   = "role-management"
)

type ResourceProvider interface {
	// Resource Management
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
