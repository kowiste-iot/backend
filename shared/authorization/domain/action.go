package domain

import (
	"backend/shared/authorization/domain/command"
	baseCmd "backend/shared/base/command"
	"context"
)

const (
	View   string = "view"
	Create string = "create"
	Update string = "update"
	Delete string = "delete"
)

type ActionProvider interface {
	CreateAction(ctx context.Context, input *command.CreateActionInput) (*Action, error)
	ListActions(ctx context.Context, input *baseCmd.BaseInput) ([]Action, error)
}

type Action struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

func AllScopes() []string {
	return []string{View, Create, Update, Delete}
}

type Actions []Action

func (s Actions) GetByName(scopesName []string) []Action {
	out := make([]Action, 0)
	for i := range s {
		for j := range scopesName {
			if s[i].Name == scopesName[j] {
				out = append(out, s[i])
			}
		}

	}
	return out
}
