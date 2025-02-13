package domain

import (
	"backend/internal/features/resource/domain/command"
	scopeDomain "backend/internal/features/scope/domain"
	"backend/pkg/config"

	baseCmd "backend/shared/base/command"
	"context"

	"github.com/google/uuid"
)

const (
	//TODO: have to keep this for now until think if is possible to use the config one.
	Measure   string = "measure-resource"
	Asset     string = "asset-resource"
	Device    string = "device-resource"
	Dashboard string = "dashboard-resource"
	Widget    string = "widget-resource"
	Action    string = "action-resource"
	Alert     string = "alert-resource"
	ResourceR string = "resource-resource"
	Tenant    string = "tenant-resource"
	Branch    string = "branch-resource"
	User      string = "user-resource"
	Role      string = "role-resource"
	Admin     string = "admin-resource"

	//
	defaultResource string = "Default Resource"
)

type ResourceProvider interface {
	CreateResource(ctx context.Context, input *baseCmd.BaseInput, resource Resource) (*Resource, error)
	GetResource(ctx context.Context, input *command.ResourceIDInput) (*Resource, error)
	ListResources(ctx context.Context, input *baseCmd.BaseInput) ([]Resource, error)
}

type Resource struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Type        string   `json:"type,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
}

func (a *Resource) SetID(id string) {
	a.ID = id
}

func New(name string, resourceType string, scopes []string, displayName string) (resource *Resource, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}

	resource = &Resource{
		ID:   id.String(),
		Name: name,
		Type: resourceType,

		Scopes:      scopes,
		DisplayName: displayName,
	}
	return
}

func NewFromRepository(id string, name string, resourceType string, scopes []string, displayName string) *Resource {
	return &Resource{
		ID:          id,
		Name:        name,
		Type:        resourceType,
		Scopes:      scopes,
		DisplayName: displayName,
	}
}
func EndpointsResources(input map[string]config.Resource) (resources []Resource) {
	for i := range input {
		scopes := scopeDomain.AllScopes()
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
