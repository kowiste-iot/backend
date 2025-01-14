package domain

import (
	"time"
)

type WSTokenClaims struct {
	TenantID  string    `json:"tid"`
	UserID    string    `json:"uid"`
	ExpiresAt time.Time `json:"exp"`
}

type TokenService interface {
	GenerateWSToken(tenantID, userID string) (string, error)
	ValidateWSToken(token string) (*WSTokenClaims, error)
}
