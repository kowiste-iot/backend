package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type TokenValidator interface {
	ValidateToken(ctx context.Context,tenant, token string) (*jwt.Token, error)
}
