package keycloak

import (
	"backend/internal/features/permission/domain"
	permissionDomain "backend/internal/features/permission/domain"
	scopeDomain "backend/internal/features/scope/domain"

	baseCmd "backend/shared/base/command"
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"

	"backend/shared/keycloak"

	"github.com/Nerzal/gocloak/v13"
)

type PermissionKeycloak struct {
	*keycloak.Keycloak
}

func New(core *keycloak.Keycloak) *PermissionKeycloak {
	return &PermissionKeycloak{
		Keycloak: core,
	}
}

func (rk PermissionKeycloak) CreatePermission(ctx context.Context, scopes []scopeDomain.Scope, input *baseCmd.BaseInput, per *domain.Permission) (*domain.Permission, error) {

	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	kcPer := newPermissionKc(per)
	created, err := createPermission(ctx, rk.Config.Host, token.AccessToken, input.TenantDomain, *input.ClientID, kcPer)

	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	// Convert back to our Permission type
	result := &domain.Permission{
		ID:               created.ID,
		Name:             created.Name,
		DecisionStrategy: created.DecisionStrategy,
		Description:      created.Description,
		Resource:         created.Resource,
		Scopes:           created.Scopes,
		Policies:         created.Policies,
	}

	return result, nil
}

func (rk PermissionKeycloak) ListPermissions(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Permission, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	kcPermissions, err := rk.Client.GetPermissions(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		gocloak.GetPermissionParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	permissions := make([]domain.Permission, len(kcPermissions))
	for i, kp := range kcPermissions {
		permissions[i] = domain.Permission{
			ID:    *kp.ID,
			Name:  *kp.Name,
			Type:  *kp.Type,
			Logic: string(*kp.Logic),
		}
		if kp.Resources != nil {
			t := *kp.Resources
			if len(t) != 1 {
				return nil, fmt.Errorf("more than 1 resource in permission")
			}
			permissions[i].Resource = t[0]
		}
		if kp.Scopes != nil {
			permissions[i].Scopes = *kp.Scopes
		}
		if kp.Resources != nil {
			permissions[i].Policies = *kp.Policies
		}
		resourPerm, err := rk.Client.GetPermissionResources(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *kp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch permission resource %w", err)
		}
		if len(resourPerm) > 0 {
			if len(resourPerm) != 1 {
				return nil, fmt.Errorf("error mor than one resource in permission")
			}
			permissions[i].Resource = *resourPerm[0].ResourceID
		}
		scopesPerm, err := rk.Client.GetPermissionScopes(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *kp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch permission scopes %w", err)
		}
		if len(scopesPerm) > 0 {
			for j := range scopesPerm {
				permissions[i].Scopes = append(permissions[i].Scopes, *scopesPerm[j].ScopeName)
			}
		}
		policyPerm, err := rk.Client.GetAuthorizationPolicyAssociatedPolicies(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *kp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch permission scopes %w", err)
		}
		if len(policyPerm) > 0 {
			for j := range policyPerm {
				permissions[i].Policies = append(permissions[i].Policies, *policyPerm[j].Name)
			}
		}
	}
	return permissions, nil
}
func (rk PermissionKeycloak) DeletePermission(ctx context.Context, input *baseCmd.BaseInput, resourceID string) (err error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, input)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	kcPermissions, err := rk.Client.GetPermissions(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		gocloak.GetPermissionParams{Resource: &resourceID},
	)
	if err != nil {
		return fmt.Errorf("failed to list permissions: %w", err)
	}

	for i := range kcPermissions {
		err = rk.Client.DeletePermission(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *kcPermissions[i].ID)
		if err != nil {
			return
		}
	}
	return
}

func createPermission(ctx context.Context, url, token, tenantID, IDofClient string, p *permissionKc) (*permissionDomain.Permission, error) {
	baseURL := url + "/admin/realms/%s/clients/%s/authz/resource-server/permission/%s"
	endpoint := fmt.Sprintf(baseURL, tenantID, IDofClient, p.Type)

	// Set default values if not provided
	if p.DecisionStrategy == "" {
		p.DecisionStrategy = permissionDomain.DecisionAffirmative
	}
	if p.Logic == "" {
		p.Logic = permissionDomain.LogicPositive
	}

	client := resty.New()

	var result permissionKc

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

	return result.ToDomain(), nil
}
