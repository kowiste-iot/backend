package keycloak

import (
	"context"
	"errors"

	tenantCmd "backend/internal/features/tenant/domain/command"
	"backend/shared/authorization/domain/command"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
)

func (k *Keycloak) ValidateToken(ctx context.Context,tenant, accessToken string) (*jwt.Token, error) {
	decodedToken, _, err := k.Client.DecodeAccessToken(
		ctx,
		accessToken,
		tenant,
	)

	if err != nil {
		return nil, err
	}

	if !decodedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return &jwt.Token{
		Raw:       decodedToken.Raw,
		Method:    decodedToken.Method,
		Header:    decodedToken.Header,
		Claims:    decodedToken.Claims,
		Signature: decodedToken.Signature,
		Valid:     decodedToken.Valid,
	}, nil
}

// ValidatePermissionUser checks if the user has permission to access a resource with specific scope
func (k *Keycloak) HasPermission(ctx context.Context, input *command.PermissionInput) (bool, error) {
	tenant := k.getTenantOrDefault(ctx)
	permissions := []string{input.Resource}
	aud := tenantCmd.ClientName(input.BranchID)
	result, err := k.Client.GetRequestingPartyPermissions(ctx,
		input.Token,
		tenant,
		gocloak.RequestingPartyTokenOptions{
			Permissions: &permissions,
			Audience:    &aud,
		})
	if err != nil {
		return false, err
	}
	hasAccess := false
	if permissions != nil {
		for _, permission := range *result {
			if permission.ResourceName != nil && *permission.ResourceName == input.Resource {
				if permission.Scopes == nil {
					// If scopes is nil means full access
					hasAccess = true
					break
				}
				for _, scope := range *permission.Scopes {
					if scope == input.Action {
						hasAccess = true
						break
					}
				}
			}
		}
	}
	return hasAccess, nil

}
