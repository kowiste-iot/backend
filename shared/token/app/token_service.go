package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"ddd/shared/token/domain"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type TokenService struct {
	secretKey []byte
	tokenTTL  time.Duration
}

func NewTokenService(secretKey string, tokenTTL time.Duration) *TokenService {
	return &TokenService{
		secretKey: []byte(secretKey),
		tokenTTL:  tokenTTL,
	}
}

func (s *TokenService) GenerateWSToken(tenantID, userID string) (string, error) {
	claims := domain.WSTokenClaims{
		TenantID:  tenantID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.tokenTTL),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, s.secretKey)
	mac.Write(payload)
	signature := mac.Sum(nil)

	token := base64.URLEncoding.EncodeToString(payload) + "." +
		base64.URLEncoding.EncodeToString(signature)

	return token, nil
}

func (s *TokenService) ValidateWSToken(token string) (*domain.WSTokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	payload, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	signature, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	// Verify signature
	mac := hmac.New(sha256.New, s.secretKey)
	mac.Write(payload)
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return nil, errors.New("invalid token signature")
	}

	var claims domain.WSTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}
