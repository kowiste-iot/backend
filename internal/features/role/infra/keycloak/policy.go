package keycloak

import (
	baseCmd "backend/shared/base/command"
	"context"
	"fmt"
	"backend/internal/features/user/domain"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-resty/resty/v2"
)



type policy struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Type             string   `json:"type"`
	Logic            string   `json:"logic"`
	DecisionStrategy string   `json:"decisionStrategy"`
	Roles            []string `json:"roles,omitempty"`
}

func (ks *RoleKeycloak) createPolicy(ctx context.Context, tenantID, clientID string, p policy) (*policy, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	roleDefs := make([]restRole, len(p.Roles))
	for i, role := range p.Roles {
		required := true
		roleDefs[i] = restRole{
			ID:       role,
			Required: required,
		}
	}

	pol := restPolicy{
		Name:       p.Name,
		Roles:      roleDefs,
		FetchRoles: true,
		Logic:      domain.LogicPositive,
	}
	createdID, err := createRolePolicy(ctx, ks.Config.Host, token.AccessToken, tenantID, clientID, pol)

	if err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	}

	return &policy{
		ID:    *createdID,
		Name:  p.Name,
		Type:  "role",
		Logic: "POSITIVE",
		Roles: p.Roles,
	}, nil
}

func (ks *RoleKeycloak) getPolicyByName(ctx context.Context, input *baseCmd.BaseInput, name string) (*policy, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.FetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	policies, err := ks.Client.GetPolicies(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, gocloak.GetPolicyParams{
		Name: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get policy by name: %w", err)
	}

	if len(policies) == 0 {
		return nil, fmt.Errorf("client not found")
	}
	roles := make([]string, 0)
	if policies[0].Roles != nil {
		for _, role := range *policies[0].Roles {
			roles = append(roles, role.String())
		}
	}

	result := policy{
		ID:               *policies[0].ID,
		Name:             *policies[0].Name,
		Description:      *policies[0].Description,
		Type:             *policies[0].Type,
		Logic:            string(*policies[0].Logic),
		DecisionStrategy: string(*policies[0].DecisionStrategy),
		Roles:            roles,
	}

	return &result, nil
}

type restRole struct {
	ID       string `json:"id"`
	Required bool   `json:"required"`
}

type restPolicy struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Roles       []restRole `json:"roles"`
	FetchRoles  bool       `json:"fetchRoles"`
	Logic       string     `json:"logic"`
}

func createRolePolicy(ctx context.Context, url, token, tenantID, clientID string, p restPolicy) (*string, error) {
	baseURL := url + "/admin/realms/%s/clients/%s/authz/resource-server/policy/role"
	endpoint := fmt.Sprintf(baseURL, tenantID, clientID)

	client := resty.New()

	var result struct {
		ID string `json:"id"`
	}

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
		return nil, fmt.Errorf("failed to create policy. Status: %d, Body: %s",
			resp.StatusCode(), string(resp.Body()))
	}

	return &result.ID, nil
}
