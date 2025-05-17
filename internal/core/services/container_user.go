package services

import (
	appUser "backend/internal/features/user/app"
	repoUser "backend/internal/features/user/infra/gorm"
	userKeycloak "backend/internal/features/user/infra/keycloak"
)

func (c *Container) initializeUserService(s *Services) error {

	userRepo := repoUser.NewRepository(c.base.DB)
	userKC := userKeycloak.New(c.auth)
	s.UserService = appUser.NewUserService(c.base, &appUser.ServiceDependencies{
		Repo: userRepo,
		Auth: userKC,
	})
	return nil
}
