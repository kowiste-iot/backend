package http

import (
	actionhandler "backend/internal/features/action/interface/rest"
	alerthandler "backend/internal/features/alert/interface/rest"
	assethandler "backend/internal/features/asset/interface/rest"
	dashboardhandler "backend/internal/features/dashboard/interface/rest"
	devicehandler "backend/internal/features/device/interface/rest"
	measurehandler "backend/internal/features/measure/interface/rest"
	rolehandler "backend/internal/features/role/interface/rest"
	tenanthandler "backend/internal/features/tenant/interface/rest"
	userhandler "backend/internal/features/user/interface/rest"
	branchhandler "backend/internal/interfaces/http/handlers/branch"
	resourcehandler "backend/internal/interfaces/http/handlers/resource"
	scopehandler "backend/internal/interfaces/http/handlers/scopes"
	"backend/pkg/config"
	"backend/shared/authentication/domain"
	"backend/shared/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	wshandler "backend/internal/interfaces/http/handlers/websocket"
)

type Server struct {
	config           *config.Config
	logger           logger.Logger
	router           *gin.Engine
	httpServer       *http.Server
	auth             domain.TokenValidator
	tenantHandler    *tenanthandler.TenantHandler
	assetHandler     *assethandler.AssetHandler
	measureHandler   *measurehandler.MeasureHandler
	dashboardHandler *dashboardhandler.DashboardHandler
	deviceHandler    *devicehandler.DeviceHandler
	actionHandler    *actionhandler.ActionHandler
	alertHandler     *alerthandler.AlertHandler
	userHandler      *userhandler.UserHandler
	rolesHandler     *rolehandler.RoleHandler
	resourceHandler  *resourcehandler.ResourceHandler
	tokenHandler     *wshandler.TokenHandler
	wsNotifyHandler  *wshandler.NotificationHandler
	scopesHandler    *scopehandler.ScopeHandler
}

type ServerDependencies struct {
	Authentication    domain.TokenValidator
	RolesHandler     *rolehandler.RoleHandler
	ResourceHandler  *resourcehandler.ResourceHandler
	BranchHandler    *branchhandler.BranchHandler
	TenantHandler    *tenanthandler.TenantHandler
	AssetHandler     *assethandler.AssetHandler
	MeasureHandler   *measurehandler.MeasureHandler
	DashboardHandler *dashboardhandler.DashboardHandler
	DeviceHandler    *devicehandler.DeviceHandler
	ActionHandler    *actionhandler.ActionHandler
	AlertHandler     *alerthandler.AlertHandler
	UserHandler      *userhandler.UserHandler
	ScopeHandler     *scopehandler.ScopeHandler
	TokenHandler     *wshandler.TokenHandler
	WSNotifyHandler  *wshandler.NotificationHandler
}

func NewServer(cfg *config.Config, logger logger.Logger, deps ServerDependencies) *Server {
	router := gin.New()
	// Setup middleware
	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		gin.Recovery(),
	)

	return &Server{
		config:           cfg,
		logger:           logger,
		router:           router,
		auth:             deps.Authentication,
		rolesHandler:     deps.RolesHandler,
		resourceHandler:  deps.ResourceHandler,
		tenantHandler:    deps.TenantHandler,
		assetHandler:     deps.AssetHandler,
		dashboardHandler: deps.DashboardHandler,
		deviceHandler:    deps.DeviceHandler,
		actionHandler:    deps.ActionHandler,
		alertHandler:     deps.AlertHandler,
		userHandler:      deps.UserHandler,
		scopesHandler:    deps.ScopeHandler,
		wsNotifyHandler:  deps.WSNotifyHandler,
		tokenHandler:     deps.TokenHandler,
		measureHandler:   deps.MeasureHandler,
		// widgetHandler:    deps.WidgetHandler,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.setupRoutes()

	addr := fmt.Sprintf("%s:%d", s.config.HTTP.Host, s.config.HTTP.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.HTTP.ReadTimeout,
		WriteTimeout: s.config.HTTP.WriteTimeout,
	}

	s.logger.Info(ctx, "Starting HTTP server on %s"+addr, map[string]interface{}{})
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
