package restkc

import (
	"backend/shared/auth/domain/permission"
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

func CreatePermission(ctx context.Context, url, token, tenantID, IDofClient string, p permission.Permission) (*permission.Permission, error) {
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

	var result permission.Permission

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
