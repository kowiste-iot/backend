package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	assethandler "ddd/internal/interfaces/http/handlers/asset"
	branchhandler "ddd/internal/interfaces/http/handlers/branch"
	resourcehandler "ddd/internal/interfaces/http/handlers/resource"
	rolehandler "ddd/internal/interfaces/http/handlers/roles"
	tenanthandler "ddd/internal/interfaces/http/handlers/tenant"
	userhandler "ddd/internal/interfaces/http/handlers/user"
	"ddd/pkg/config"
	auth "ddd/shared/auth/domain"
	"ddd/shared/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	wshandler "ddd/internal/interfaces/http/handlers/websocket"
)

type Server struct {
	config          *config.Config
	logger          logger.Logger
	router          *gin.Engine
	httpServer      *http.Server
	auth            auth.AuthProvider
	tenantHandler   *tenanthandler.TenantHandler
	assetHandler    *assethandler.AssetHandler
	userHandler     *userhandler.UserHandler
	rolesHandler    *rolehandler.RoleHandler
	resourceHandler *resourcehandler.ResourceHandler
	tokenHandler    *wshandler.TokenHandler
	wsNotifyHandler *wshandler.NotificationHandler
}

type ServerDependencies struct {
	Authentication  auth.AuthProvider
	RolesHandler    *rolehandler.RoleHandler
	ResourceHandler *resourcehandler.ResourceHandler
	BranchHandler   *branchhandler.BranchHandler
	TenantHandler   *tenanthandler.TenantHandler
	AssetHandler    *assethandler.AssetHandler
	UserHandler     *userhandler.UserHandler
	TokenHandler    *wshandler.TokenHandler
	WSNotifyHandler *wshandler.NotificationHandler
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
		config:          cfg,
		logger:          logger,
		router:          router,
		auth:            deps.Authentication,
		rolesHandler:    deps.RolesHandler,
		resourceHandler: deps.ResourceHandler,
		tenantHandler:   deps.TenantHandler,
		assetHandler:    deps.AssetHandler,
		userHandler:     deps.UserHandler,
		wsNotifyHandler: deps.WSNotifyHandler,
		tokenHandler:    deps.TokenHandler,
		// measureHandler:   deps.MeasureHandler,
		// dashboardHandler: deps.DashboardHandler,
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
