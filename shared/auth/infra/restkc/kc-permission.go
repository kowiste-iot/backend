package restkc

import (
	"context"
	"ddd/shared/auth/domain/permission"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Permission struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Type             string   `json:"type"`
	Resources        []string `json:"resources"`
	Scopes           []string `json:"scopes"`
	Policies         []string `json:"policies"`
	DecisionStrategy string   `json:"decisionStrategy"`
	Logic            string   `json:"logic"`
}

func CreatePermission(ctx context.Context, url, token, tenantID, IDofClient string, p Permission) (*Permission, error) {
	baseURL := url + "/admin/realms/%s/clients/%s/authz/resource-server/permission/%s"
	endpoint := fmt.Sprintf(baseURL, tenantID, IDofClient, p.Type)

	// Set default values if not provided
	if p.DecisionStrategy == "" {
		p.DecisionStrategy = permission.DecisionAffirmative
	}
	if p.Logic == "" {
		p.Logic = permission.LogicPositive
	}

	client := resty.New()

	var result Permission

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(p).
		SetResult(&result).
		Post(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != 201 {
		return nil, fmt.Errorf("failed to create permission. Status: %d, Body: %s",
			resp.StatusCode(), string(resp.Body()))
	}

	return &result, nil
}
