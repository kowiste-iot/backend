package restkc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type ClientSettings struct {
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

// {
//     "id": "8baf6327-08e4-4961-afd1-d3c1edd66eb1",
//     "clientId": "8baf6327-08e4-4961-afd1-d3c1edd66eb1",
//     "name": "main-service",
//     "allowRemoteResourceManagement": true,
//     "policyEnforcementMode": "ENFORCING",
//     "resources": [],
//     "policies": [],
//     "scopes": [],
//     "decisionStrategy": "AFFIRMATIVE"
// }
//"http://localhost:7080/auth/auth/admin/realms/%s/clients/%s/authz/resource-server"

// {
//   "id": "7de4f407-1509-4f99-92d3-45bdf23d39c5",
//   "clientId": "7de4f407-1509-4f99-92d3-45bdf23d39c5",
//   "name": "main-service",
//   "allowRemoteResourceManagement": true,
//   "policyEnforcementMode": "ENFORCING",
//   "resources": [],
//   "policies": [],
//   "scopes": [],
//   "decisionStrategy": "AFFIRMATIVE"
// }
func UpdateClientSettings(ctx context.Context, url, token, tenantID, IDofClient string, p *ClientSettings) error {
	baseURL := url + "/admin/realms/%s/clients/%s/authz/resource-server"
	endpoint := fmt.Sprintf(baseURL, tenantID, IDofClient)

	client := resty.New()
	b, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
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
