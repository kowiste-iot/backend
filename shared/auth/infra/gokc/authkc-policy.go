package keycloak

import (
	"backend/shared/auth/domain/command"
	"backend/shared/auth/domain/permission"
	"backend/shared/auth/domain/policy"
	"backend/shared/auth/infra/restkc"
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

func (ks *KeycloakService) CreatePolicy(ctx context.Context, tenantID, clientID string, p policy.Policy) (*policy.Policy, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	roleDefs := make([]restkc.Role, len(p.Roles))
	for i, role := range p.Roles {
		required := true
		roleDefs[i] = restkc.Role{
			ID:       role,
			Required: required,
		}
	}

	pol := restkc.Policy{
		Name:       p.Name,
		Roles:      roleDefs,
		FetchRoles: true,
		Logic:      permission.LogicPositive,
	}
	createdID, err := restkc.CreateRolePolicy(ctx, ks.config.Host, token.AccessToken, tenantID, clientID, pol)

	if err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	}

	return &policy.Policy{
		ID:    *createdID,
		Name:  p.Name,
		Type:  "role",
		Logic: "POSITIVE",
		Roles: p.Roles,
	}, nil
}

func (ks *KeycloakService) UpdatePolicy(ctx context.Context, tenantID, clientID string, p policy.Policy) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	roleDefs := make([]gocloak.RoleDefinition, len(p.Roles))
	for i, role := range p.Roles {
		required := true
		roleDefs[i] = gocloak.RoleDefinition{
			ID:       &role,
			Required: &required,
		}
	}

	policyType := policy.TypeRole
	logic := gocloak.POSITIVE
	kcPolicy := gocloak.PolicyRepresentation{
		ID:    &p.ID,
		Type:  &policyType,
		Logic: logic,
		Name:  &p.Name,
		RolePolicyRepresentation: gocloak.RolePolicyRepresentation{
			Roles: &roleDefs,
		},
	}

	err = ks.client.UpdatePolicy(ctx, token.AccessToken, tenantID, clientID, kcPolicy)
	if err != nil {
		return fmt.Errorf("failed to update policy: %w", err)
	}

	return nil
}
func (ks *KeycloakService) DeletePolicy(ctx context.Context, input *command.PolicyNameInput) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.fetchClient(ctx, &input.BaseInput)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	policy, err := ks.GetPolicyByName(ctx, input)
	if err != nil {
		return fmt.Errorf("error fetching policy: %w", err)
	}
	err = ks.client.DeletePolicy(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, policy.ID)
	if err != nil {
		return fmt.Errorf("failed to update policy: %w", err)
	}

	return nil
}

func (ks *KeycloakService) GetPolicyByName(ctx context.Context, input *command.PolicyNameInput) (*policy.Policy, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.fetchClient(ctx, &input.BaseInput)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	policies, err := ks.client.GetPolicies(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, gocloak.GetPolicyParams{
		Name: &input.PolicyName,
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

	result := policy.Policy{
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
