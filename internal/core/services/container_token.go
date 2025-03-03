package services

import (
	"backend/shared/token/app"
	"backend/shared/token/domain"
	"backend/shared/token/infra/keycloak"
)

func (c *Container) initializeTokenService(s *Services) (err error) {
	tokenKc := keycloak.New(&domain.TokenConfiguration{}, c.auth)
	s.TokenService = app.New(c.base, tokenKc)
	return
}
