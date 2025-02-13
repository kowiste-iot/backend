package domain

import (
	baseCmd "backend/shared/base/command"
	"context"

	"github.com/google/uuid"
)

const (
	View   string = "view"
	Create string = "create"
	Update string = "update"
	Delete string = "delete"
)

type ScopeProvider interface {
	CreateScope(ctx context.Context, input *baseCmd.BaseInput, scope Scope) (*Scope, error)
	ListScopes(ctx context.Context, input *baseCmd.BaseInput) ([]Scope, error)
}

type Scope struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

func New(name string, displayName string) (scope *Scope, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}

	scope = &Scope{
		ID:          id.String(),
		Name:        name,
		DisplayName: displayName,
	}
	return
}

func NewFromRepository(id string, name string, displayName string) *Scope {
	return &Scope{
		ID:          id,
		Name:        name,
		DisplayName: displayName,
	}
}
func AllScopes() []string {
	return []string{View, Create, Update, Delete}
}
