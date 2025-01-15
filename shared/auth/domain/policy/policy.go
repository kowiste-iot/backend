package policy

import "context"

const (
	TypeRole  string = "role"
	Enforcing string = "ENFORCING"
)

type Policy struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Type             string   `json:"type"`
	Logic            string   `json:"logic"`
	DecisionStrategy string   `json:"decisionStrategy"`
	Roles            []string `json:"roles,omitempty"`
}

type PolicyProvider interface {
	CreatePolicy(ctx context.Context, tenantID, clientID string, policy Policy) (*Policy, error)
	UpdatePolicy(ctx context.Context, tenantID, clientID string, policy Policy) error
}
