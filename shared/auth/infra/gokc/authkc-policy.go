package keycloak

import (
	"context"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/policy"
	"ddd/shared/auth/infra/restkc"
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

// func (ks *KeycloakService) CreatePolicy(ctx context.Context, tenantID, clientID string, p policy.Policy) (*policy.Policy, error) {
// 	token, err := ks.GetValidToken(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get token: %w", err)
// 	}

// 	roleDefs := make([]gocloak.RoleDefinition, len(p.Roles))
// 	for i, role := range p.Roles {
// 		required := false
// 		roleDefs[i] = gocloak.RoleDefinition{
// 			ID:       &role,
// 			Required: &required,
// 		}
// 	}

// 	policyType := "role"
// 	logic := gocloak.POSITIVE
// 	//TODO: add fetchRoles to this struct
// 	configPolicy := map[string]string{
// 		"fetchRoles": "",
// 	}
// 	kcPolicy := gocloak.PolicyRepresentation{
// 		Config: &configPolicy,
// 		Type:   &policyType,
// 		Logic:  logic,
// 		Name:   &p.Name,
// 		RolePolicyRepresentation: gocloak.RolePolicyRepresentation{
// 			Roles: &roleDefs,
// 		},
// 	}

// 	createdPolicy, err := ks.client.CreatePolicy(
// 		ctx,
// 		token.AccessToken,
// 		tenantID,
// 		clientID,
// 		kcPolicy,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create policy: %w", err)
// 	}

// 	return &policy.Policy{
// 		ID:    *createdPolicy.ID,
// 		Name:  *createdPolicy.Name,
// 		Type:  *createdPolicy.Type,
// 		Logic: string(*createdPolicy.Logic),
// 		Roles: p.Roles,
// 	}, nil
// }

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
