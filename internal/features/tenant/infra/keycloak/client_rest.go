package keycloak

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type clientSettings struct {
	ID                            string   `json:"id"`
	ClientID                      string   `json:"clientId"`
	Name                          string   `json:"name"`
	AllowRemoteResourceManagement bool     `json:"allowRemoteResourceManagement"`
	PolicyEnforcementMode         string   `json:"policyEnforcementMode"`
	Resources                     []string `json:"resources"`
	Policies                      []string `json:"policies"`
	Scopes                        []string `json:"scopes"`
	DecisionStrategy              string   `json:"decisionStrategy"`
}

func (ks *BranchKeycloak) updateClientSettings(ctx context.Context, url, token, tenantID, IDofClient string, p *clientSettings) error {
	baseURL := url + "/admin/realms/%s/clients/%s/authz/resource-server"
	endpoint := fmt.Sprintf(baseURL, tenantID, IDofClient)

	client := resty.New()
	// No need to set result since the response is empty (204 No Content)
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(p).
		Put(endpoint)

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("failed to update client settings. Status: %d, Body: %s",
			resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
