package http

import (
	"backend/pkg/config"
	"backend/shared/authentication/domain"
	"backend/shared/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     *config.Config
	logger     logger.Logger
	router     *gin.Engine
	httpServer *http.Server
	auth       domain.TokenValidator
}

type ServerDependencies struct {
	Authentication domain.TokenValidator
}

func NewServer(cfg *config.Config, logger logger.Logger, deps ServerDependencies) *Server {
	router := gin.New()
	// Setup middleware
	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     cfg.HTTP.CORSAllowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Tenant-ID", "X-Branch-ID"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		gin.Recovery(),
	)

	return &Server{
		config: cfg,
		logger: logger,
		router: router,
		auth:   deps.Authentication,
	}
}

func (s *Server) Start(ctx context.Context) error {

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

func (s *Server) GetRouter(ctx context.Context) *gin.Engine {
	return s.router
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
