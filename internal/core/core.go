package core

import (
	"backend/pkg/config"
	"backend/shared/base"
	"backend/shared/logger"
	"backend/shared/logger/openob"

	"backend/shared/validator"
	"context"
	"errors"
	"fmt"
	"time"

	streamDomain "backend/shared/stream/domain"
	"backend/shared/stream/infra/nats"
	repoStream "backend/shared/stream/repo"

	appAsset "backend/internal/features/asset/app"
	repoAsset "backend/internal/features/asset/infra/gorm"
	assethandler "backend/internal/features/asset/interface/rest"

	appDashboard "backend/internal/features/dashboard/app"
	repoDashboard "backend/internal/features/dashboard/infra/gorm"
	dashboardhandler "backend/internal/features/dashboard/interface/rest"

	appDevice "backend/internal/features/device/app"
	repoDevice "backend/internal/features/device/infra/gorm"
	devicehandler "backend/internal/features/device/interface/rest"

	appAction "backend/internal/features/action/app"
	repoAction "backend/internal/features/action/infra/gorm"
	actionhandler "backend/internal/features/action/interface/rest"

	appAlert "backend/internal/features/alert/app"
	repoAlert "backend/internal/features/alert/infra/gorm"
	alerthandler "backend/internal/features/alert/interface/rest"

	appMeasure "backend/internal/features/measure/app"
	repoMeasure "backend/internal/features/measure/infra/gorm"
	measurehandler "backend/internal/features/measure/interface/rest"

	appTenant "backend/internal/features/tenant/app"
	repoTenant "backend/internal/features/tenant/infra/gorm"
	tenantKeycloak "backend/internal/features/tenant/infra/keycloak"
	tenanthandler "backend/internal/features/tenant/interface/rest"

	appUser "backend/internal/features/user/app"
	repoUser "backend/internal/features/user/infra/gorm"
	userKeycloak "backend/internal/features/user/infra/keycloak"
	userhandler "backend/internal/features/user/interface/rest"

	appRole "backend/internal/features/role/app"
	roleKeycloak "backend/internal/features/role/infra/keycloak"
	rolehandler "backend/internal/features/role/interface/rest"

	appScope "backend/internal/features/scope/app"
	scopeKeycloak "backend/internal/features/scope/infra/keycloak"
	scopehandler "backend/internal/features/scope/interface/rest"

	appResource "backend/internal/features/resource/app"
	resourceKeycloak "backend/internal/features/resource/infra/keycloak"
	resourcehandler "backend/internal/features/resource/interface/rest"

	appPermission "backend/internal/features/permission/app"
	permissionKeycloak "backend/internal/features/permission/infra/keycloak"

	"backend/internal/interfaces/http"

	// wshandler "backend/internal/interfaces/http/handlers/websocket"
	// appToken "backend/shared/token/app"
	// appWS "backend/shared/websocket/app"
	kcCore "backend/shared/keycloak"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Core struct {
	cfg    *config.Config
	logger logger.Logger
	db     *gorm.DB
	server *http.Server
	stream streamDomain.StreamClient
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

// In initServer or create a new init function
func (c *Core) initStreaming() error {
	msgRepo := repoStream.NewMessageRepository(c.db)
	streamConfig := &streamDomain.StreamConfig{
		URL:            "http://localhost:4222", // Add this to your config
		MaxReconnects:  5,
		ReconnectWait:  time.Second * 5,
		ConnectTimeout: time.Second * 2,
		WriteTimeout:   time.Second * 2,
		PersistMessage: true,
	}

	factory := nats.NewNatsClientFactory(msgRepo)
	streamClient, err := factory.CreateClient(streamConfig)
	if err != nil {
		return fmt.Errorf("failed to create stream client: %w", err)
	}

	c.stream = streamClient
	return nil
}

func (c *Core) initServer(ctx context.Context) error {

	//load authentication config
	tenantConfig, err := config.LoadTenant()
	if err != nil {
		return errors.New("cant load tenant config")
	}
	validator.InitValidator(tenantConfig.Authorization.Roles)

	kCore, err := kcCore.New(&kcCore.KeycloakConfig{
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
		Perm:   kCore,
	}

	//Asset
	assetRepo := repoAsset.NewRepository(c.db)
	assetService := appAsset.NewService(base, assetRepo)

	assetDepRepo := repoAsset.NewDependencyRepository(c.db)
	assetDep := appAsset.NewAssetDependencyService(base, assetDepRepo)
	//Measure
	measureRepo := repoMeasure.NewRepository(c.db)
	measureService := appMeasure.NewService(base, measureRepo, assetDep)
	//Dashboard
	dashboardRepo := repoDashboard.NewRepository(c.db)
	dashboardService := appDashboard.NewService(base, dashboardRepo, assetDep)
	//Widget
	widgetRepo := repoDashboard.NewWidgetRepository(c.db)
	widgetService := appDashboard.NewWidgetService(base, widgetRepo)
	//Device
	deviceRepo := repoDevice.NewRepository(c.db)
	deviceService := appDevice.NewService(base, deviceRepo, assetDep)
	//Action
	actionRepo := repoAction.NewRepository(c.db)
	actionService := appAction.NewService(base, actionRepo, assetDep)
	//Alert
	alertRepo := repoAlert.NewRepository(c.db)
	alertService := appAlert.NewService(base, alertRepo, assetDep)
	//User
	userRepo := repoUser.NewRepository(c.db)
	userKC := userKeycloak.New(kCore)
	userService := appUser.NewService(base, &appUser.ServiceDependencies{
		Repo: userRepo,
		Auth: userKC,
	})
	//Roles
	roleKC := roleKeycloak.New(kCore)
	roleService := appRole.NewService(base, roleKC, appRole.Config{
		DefaultRoles: tenantConfig.Authorization.Roles,
	})
	//Roles
	scopeKC := scopeKeycloak.New(kCore)
	scopeService := appScope.NewService(base, &appScope.ServiceDependencies{
		Repo: scopeKC,
	})

	//Permission
	permissionKC := permissionKeycloak.New(kCore)
	permissionService := appPermission.NewService(base, &appPermission.ServiceDependencies{
		Repo:  permissionKC,
		Scope: scopeService,
		Role:  roleService,
		Config: &appPermission.Config{
			DefaultRoles: tenantConfig.Authorization.Roles,
		},
	})
	//Resource
	resourceKC := resourceKeycloak.New(kCore)
	resourceService := appResource.NewService(base, &appResource.ServiceDependencies{
		Repo:       resourceKC,
		Roles:      roleService,
		Permission: permissionService,
		Scopes:     scopeService,
		Config: &appResource.Config{
			DefaultRoles: tenantConfig.Authorization.Roles,
		},
	})
	//Branch
	branchRepo := repoTenant.NewBranchRepository(c.db)
	branchKC := tenantKeycloak.NewBranch(tenantConfig, kCore)

	branchService := appTenant.NewBranchService(base, &appTenant.BranchDependencies{
		Branch:     branchKC,
		Repo:       branchRepo,
		Role:       roleService,
		Scope:      scopeService,
		Resource:   resourceService,
		Permission: permissionService,
		Config:     tenantConfig,
	})
	//Tenant

	tenantRepo := repoTenant.NewTenantRepository(c.db)
	tenantKC := tenantKeycloak.New(tenantConfig, kCore)
	tenantDep := appTenant.ServiceDependencies{
		Branch: branchService,
		Tenant: tenantKC,
		Repo:   tenantRepo,
		User:   userService,
	}
	tenantService := appTenant.NewTenantService(base, &tenantDep)

	//Websocket
	// appT := appToken.NewTokenService("wA7pH9#kL$mN4@vX2*qR8", 8*time.Hour)
	// appH := appWS.NewHub()
	deps := http.ServerDependencies{
		Authentication: kCore,
		BranchHandler: tenanthandler.NewBranch(tenanthandler.BranchDependencies{
			Logger:        c.logger,
			BranchService: branchService,
		}),
		TenantHandler: tenanthandler.New(tenanthandler.Dependencies{
			Logger:        c.logger,
			TenantService: tenantService,
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
		WidgetHandler: dashboardhandler.NewWidget(dashboardhandler.DependenciesWidget{
			Logger:        c.logger,
			WidgetService: widgetService,
		}),
		DeviceHandler: devicehandler.New(devicehandler.Dependencies{
			Logger:        c.logger,
			DeviceService: deviceService,
		}),
		ActionHandler: actionhandler.New(actionhandler.Dependencies{
			Logger:        c.logger,
			ActionService: actionService,
		}),
		AlertHandler: alerthandler.New(alerthandler.Dependencies{
			Logger:       c.logger,
			AlertService: alertService,
		}),
		UserHandler: userhandler.New(userhandler.Dependencies{
			Logger:      c.logger,
			UserService: userService,
		}),
		RolesHandler: rolehandler.New(rolehandler.Dependencies{
			Logger:      c.logger,
			RoleService: roleService,
		}),
		ResourceHandler: resourcehandler.New(resourcehandler.Dependencies{
			Logger:   c.logger,
			Resource: resourceService,
		}),
		ScopeHandler: scopehandler.New(scopehandler.Dependencies{
			Logger:       c.logger,
			ScopeService: scopeService,
		}),
		// TokenHandler:    wshandler.NewTokenHandler(appT),
		// WSNotifyHandler: wshandler.NewNotificationHandler(appH, natsClient, appT),
	}

	c.server = http.NewServer(c.cfg, c.logger, deps)
	return nil
}

func (c *Core) Start(ctx context.Context) error {
	return c.server.Start(ctx)
}
