package keycloak

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type ProtocolMapperConfig struct {
	ClaimName               string `json:"claim.name"`
	FullPath                string `json:"full.path"`
	IDTokenClaim            string `json:"id.token.claim"`
	AccessTokenClaim        string `json:"access.token.claim"`
	LightweightClaim        string `json:"lightweight.claim"`
	UserinfoTokenClaim      string `json:"userinfo.token.claim"`
	IntrospectionTokenClaim string `json:"introspection.token.claim"`
}

type ProtocolMapper struct {
	Protocol       string               `json:"protocol"`
	ProtocolMapper string               `json:"protocolMapper"`
	Name           string               `json:"name"`
	Config         ProtocolMapperConfig `json:"config"`
}

func CreateProtocolMapper(ctx context.Context, url, token, realmName, clientID string, mapper ProtocolMapper) error {
	baseURL := fmt.Sprintf("%s/admin/realms/%s/clients/%s/protocol-mappers/models",
		strings.TrimRight(url, "/"), realmName, clientID)

	client := resty.New().
		SetTimeout(30 * time.Second)

	resp, err := client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(mapper).
		Post(baseURL)

	if err != nil {
		return fmt.Errorf("creating protocol mapper: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
