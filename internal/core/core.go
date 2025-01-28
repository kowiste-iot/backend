package core

import (
	"context"
	appAsset "ddd/internal/features/asset/app"
	repoAsset "ddd/internal/features/asset/infra/gorm"
	assethandler "ddd/internal/features/asset/interface/rest"

	appDashboard "ddd/internal/features/dashboard/app"
	repoDashboard "ddd/internal/features/dashboard/infra/gorm"
	dashboardhandler "ddd/internal/features/dashboard/interface/rest"

	appMeasure "ddd/internal/features/measure/app"
	repoMeasure "ddd/internal/features/measure/infra/gorm"
	measurehandler "ddd/internal/features/measure/interface/rest"
	"errors"

	appTenant "ddd/internal/features/tenant/app"
	repoTenant "ddd/internal/features/tenant/infra/gorm"
	tenanthandler "ddd/internal/features/tenant/interface/rest"

	appUser "ddd/internal/features/user/app"
	repoUser "ddd/internal/features/user/infra/gorm"
	userhandler "ddd/internal/features/user/interface/rest"

	"ddd/internal/interfaces/http"
	branchhandler "ddd/internal/interfaces/http/handlers/branch"
	resourcehandler "ddd/internal/interfaces/http/handlers/resource"
	rolehandler "ddd/internal/interfaces/http/handlers/roles"
	scopehandler "ddd/internal/interfaces/http/handlers/scopes"
	wshandler "ddd/internal/interfaces/http/handlers/websocket"
	"ddd/pkg/config"
	appAuth "ddd/shared/auth/app"
	keycloak "ddd/shared/auth/infra/gokc"
	"ddd/shared/base"
	"ddd/shared/logger"
	"ddd/shared/logger/openob"
	"ddd/shared/streaming/domain"
	"ddd/shared/streaming/infrastructure/nats"
	"ddd/shared/validator"

	//appNats "ddd/shared/nats/app"
	appToken "ddd/shared/token/app"
	appWS "ddd/shared/websocket/app"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Core struct {
	cfg    *config.Config
	logger logger.Logger
	db     *gorm.DB
	server *http.Server
}

func NewCore(ctx context.Context) (*Core, error) {

	core := &Core{}
	if err := core.initConfig(); err != nil {
		return nil, err
	}

	if err := core.initLogger(ctx); err != nil {
		return nil, err
	}

	if err := core.initDB(ctx); err != nil {
		return nil, err
	}

	if err := core.initServer(ctx); err != nil {
		return nil, err
	}

	return core, nil
}

func (c *Core) initConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c.cfg = cfg
	return nil
}

func (c *Core) initLogger(ctx context.Context) error {
	logger, err := openob.NewLogger(openob.Config{
		ServiceName:   c.cfg.App.Name,
		Environment:   c.cfg.App.Environment,
		Endpoint:      c.cfg.Telemetry.Endpoint,
		Headers:       c.cfg.Telemetry.Headers,
		TenantID:      c.cfg.Telemetry.TenantID,
		StreamName:    c.cfg.Telemetry.StreamName,
		ConsoleOutput: true,
		EnableTracing: c.cfg.Telemetry.TracingEnabled,
	})
	if err != nil {
		return err
	}
	c.logger = logger
	return nil
}

func (c *Core) initDB(ctx context.Context) error {
	db, err := gorm.Open(sqlite.Open("iot.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

func (c *Core) initServer(ctx context.Context) error {

	//load authentication config
	tenantConfig, err := config.LoadTenant()
	if err != nil {
		return errors.New("cant load tenant config")
	}
	validator.InitValidator(tenantConfig.Authorization.Roles)

	kc, err := keycloak.NewKeycloakService(keycloak.KeycloakConfig{
		Host:         c.cfg.Authentication.Host,
		Realm:        c.cfg.Authentication.Realm,
		ClientID:     c.cfg.Authentication.ClientID,
		ClientSecret: c.cfg.Authentication.ClientSecret,
		WebClientID:  c.cfg.Authentication.WebClientID,
	})
	if err != nil {
		return err
	}

	base := &base.BaseService{
		Logger: c.logger,
		Auth:   kc,
	}

	authService := appAuth.NewAuthService(tenantConfig, base, kc, kc, kc, kc, kc, kc, kc, kc)

	//Branch
	branchRepo := repoTenant.NewBranchRepository(c.db)
	branchService := appTenant.NewBranchService(base, authService, branchRepo)

	//Asset
	assetRepo := repoAsset.NewRepository(c.db)
	assetService := appAsset.NewService(base, assetRepo)
	//Measure
	measureRepo := repoMeasure.NewRepository(c.db)
	measureService := appMeasure.NewService(base, measureRepo)
	//Dashboard
	dashboardRepo := repoDashboard.NewRepository(c.db)
	dashboardService := appDashboard.NewService(base, dashboardRepo)
	//User
	userRepo := repoUser.NewRepository(c.db)
	userService := appUser.NewService(base, kc, userRepo)

	//Roles

	//Tenant

	tenantRepo := repoTenant.NewTenantRepository(c.db)
	tenantDep := appTenant.ServiceDependencies{
		Branch: branchService,
		Auth:   authService,
		Repo:   tenantRepo,
		User:   userService,
	}
	tenantService := appTenant.NewTenantService(base, &tenantDep)

	//Nats
	conn := nats.NewConnection(domain.ConnectionConfig{
		URL: "http://localhost:4222",
	})
	conn.Connect()
	natsClient := nats.NewClient(conn)
	//Websocket
	appT := appToken.NewTokenService("wA7pH9#kL$mN4@vX2*qR8", 8*time.Hour)
	appH := appWS.NewHub()
	deps := http.ServerDependencies{
		Authentication: kc,
		BranchHandler: branchhandler.New(branchhandler.Dependencies{
			Logger:        c.logger,
			BranchService: branchService,
		}),
		TenantHandler: tenanthandler.New(tenanthandler.Dependencies{
			Logger:        c.logger,
			TenantService: tenantService,
			AuthService:   authService,
		}),
		AssetHandler: assethandler.New(assethandler.Dependencies{
			Logger:       c.logger,
			AssetService: assetService,
		}),
		MeasureHandler: measurehandler.New(measurehandler.Dependencies{
			Logger:         c.logger,
			MeasureService: measureService,
		}),
		DashboardHandler: dashboardhandler.New(dashboardhandler.Dependencies{
			Logger:           c.logger,
			DashboardService: dashboardService,
		}),
		UserHandler: userhandler.New(userhandler.Dependencies{
			Logger:      c.logger,
			UserService: userService,
		}),
		RolesHandler: rolehandler.New(rolehandler.Dependencies{
			Logger:      c.logger,
			AuthService: authService,
		}),
		ResourceHandler: resourcehandler.New(resourcehandler.Dependencies{
			Logger:      c.logger,
			AuthService: authService,
		}),
		ScopeHandler: scopehandler.New(scopehandler.Dependencies{
			Logger:      c.logger,
			AuthService: authService,
		}),

		TokenHandler:    wshandler.NewTokenHandler(appT),
		WSNotifyHandler: wshandler.NewNotificationHandler(appH, natsClient, appT),
	}

	c.server = http.NewServer(c.cfg, c.logger, deps)
	return nil
}

func (c *Core) Start(ctx context.Context) error {
	return c.server.Start(ctx)
}
