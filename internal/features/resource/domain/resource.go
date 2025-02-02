package domain

import (
	"backend/pkg/config"
	"backend/shared/auth/domain/scope"
	resourceCmd "backend/shared/authorization/domain/command"
	baseCmd "backend/shared/base/command"
	"context"
)

const (
	//TODO: have to keep this for now until think if is possible to use the config one.
	Measure   string = "measure-resource"
	Asset     string = "asset-resource"
	Device    string = "device-resource"
	Dashboard string = "dashboard-resource"
	Action    string = "action-resource"
	Alert     string = "alert-resource"
	Tenant    string = "tenant-resource"
	Branch    string = "branch-resource"
	User      string = "user-resource"
	Role      string = "role-resource"
	Admin     string = "admin-resource"

	//
	defaultResource string = "Default Resource"
)

type ResourceProvider interface {
	CreateResource(ctx context.Context, tenantID, clientID string, resource Resource) (*Resource, error) //Should we allow crate resource?
	GetResource(ctx context.Context, input *resourceCmd.ResourceIDInput) (*Resource, error)
	ListResources(ctx context.Context, input *baseCmd.BaseInput) ([]Resource, error)
	AssignRoleToResource(ctx context.Context, input *resourceCmd.ResourceAssignRoleInput) error
	RemoveRolesFromResource(ctx context.Context, input *resourceCmd.ResourceAssignRoleInput) error
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

func EndpointsResources(input map[string]config.Resource) (resources []Resource) {
	for i := range input {
		scopes := scope.AllScopes()
		if input[i].Scopes != nil {
			scopes = *input[i].Scopes //This allow resources to have limit scopes for ex only view and create
		}
		resources = append(resources, Resource{
			Name:   input[i].Name,
			Type:   input[i].Type,
			Scopes: scopes,
		})
	}
	return
}
