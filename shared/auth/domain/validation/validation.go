package validation

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type AuthProvider interface {
	ValidateToken(ctx context.Context, token string) (*jwt.Token, error)
	ValidatePermissionService(ctx context.Context, token, clientID, resource, scope string) (bool, error)
	ValidatePermissionUser(ctx context.Context, token, clientID, resource, scope string) (bool, error)
}
