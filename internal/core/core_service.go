package core

import (
	"backend/pkg/config"
	"backend/shared/base"
	websocket_handler "backend/shared/websocket/interface/rest"
	"fmt"

	"backend/shared/validator"
	"context"
	"errors"

	"backend/internal/core/services"
	actionhandler "backend/internal/features/action/interface/rest"
	alerthandler "backend/internal/features/alert/interface/rest"
	assethandler "backend/internal/features/asset/interface/rest"
	dashboardhandler "backend/internal/features/dashboard/interface/rest"
	devicehandler "backend/internal/features/device/interface/rest"
	ingesthandler "backend/internal/features/ingest/interface/rest"
	measurehandler "backend/internal/features/measure/interface/rest"
	resourcehandler "backend/internal/features/resource/interface/rest"
	rolehandler "backend/internal/features/role/interface/rest"
	scopehandler "backend/internal/features/scope/interface/rest"
	userhandler "backend/internal/features/user/interface/rest"

	tenanthandler "backend/internal/features/tenant/interface/rest"

	"backend/internal/interfaces/http"
	"backend/internal/interfaces/http/middleware"
)

func (c *Core) initServer(ctx context.Context) error {

	//load authentication config
	tenantConfig, err := config.LoadTenant()
	if err != nil {
		return errors.New("cant load tenant config")
	}
	validator.InitValidator(tenantConfig.Authorization.Roles)

	base := &base.BaseService{
		Logger: c.logger,
		DB:     c.db,
		Perm:   c.auth,
	}
	//init services
	serviceContainer := services.NewContainer(base, c.auth, tenantConfig, c.cfg)
	services, err := serviceContainer.Initialize()
	if err != nil {
		return err
	}

	c.server = http.NewServer(c.cfg, c.logger, http.ServerDependencies{
		Authentication: c.auth,
	})

	middleware := middleware.NewMiddlewareManager(c.logger, c.auth, services.TenantService)

	//Init routes
	route := c.server.GetRouter(ctx)
	apiRoute := route.Group("api")
	apiRoute.Use(middleware.Auth())

	tenatHandler := tenanthandler.New(base, services.TenantService, middleware)
	apiTenant := tenatHandler.Init(apiRoute)
	branchHandler := tenanthandler.NewBranch(base, services.BranchService, middleware)
	apiBranch := branchHandler.Init(apiTenant)
	fmt.Println(apiBranch)
	//Asset
	assetHandler := assethandler.New(base, services.AssetService)
	assetHandler.Init(apiRoute)
	//Measure
	mhandler := measurehandler.New(base, services.MeasureService)
	mhandler.Init(apiRoute)
	//Dashboard
	dashhandler := dashboardhandler.New(base, services.DashboardService)
	dashboardRoute := dashhandler.Init(apiRoute)
	widHandler := dashboardhandler.NewWidget(base, services.WidgetService)
	widHandler.Init(dashboardRoute)
	//Device
	devHandler := devicehandler.New(base, services.DeviceService)
	devHandler.Init(apiRoute)
	//Action
	actHandler := actionhandler.New(base, services.ActionService)
	actHandler.Init(apiRoute)
	//Alert
	alertHandler := alerthandler.New(base, services.AlertService)
	alertHandler.Init(apiRoute)
	//User
	userHandler := userhandler.New(base, services.UserService)
	userHandler.Init(apiRoute)
	//Role
	roleHandler := rolehandler.New(base, services.RoleService)
	roleHandler.Init(apiRoute)
	//Scope
	scopeHandler := scopehandler.New(base, services.ScopeService)
	scopeHandler.Init(apiRoute)
	//Resource
	resourceHandler := resourcehandler.New(base, services.ResourceService)
	resourceHandler.Init(apiRoute)
	//Ingest
	ingestHandler := ingesthandler.New(base, services.IngestService)
	ingestHandler.Init(apiRoute)
	//Websocket
	wsHandler := websocket_handler.New(base, services.WebSocketHub, services.WebSocketService, services.TokenService)
	wsHandler.Init(apiRoute)
	return nil
}
