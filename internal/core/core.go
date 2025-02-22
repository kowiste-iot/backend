package core

import (
	"backend/pkg/config"
	"backend/shared/logger"
	"backend/shared/logger/openob"

	"context"
	"fmt"
	"time"

	streamDomain "backend/shared/stream/domain"
	"backend/shared/stream/infra/nats"
	repoStream "backend/shared/stream/repo"

	"backend/internal/interfaces/http"
	kcCore "backend/shared/keycloak"

	// wshandler "backend/internal/interfaces/http/handlers/websocket"
	// appToken "backend/shared/token/app"
	// appWS "backend/shared/websocket/app"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Core struct {
	cfg    *config.Config
	logger logger.Logger
	db     *gorm.DB
	auth   *kcCore.Keycloak
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
	if err := core.initAuth(); err != nil {
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

func (c *Core) initAuth() (err error) {
	c.auth, err = kcCore.New(&kcCore.KeycloakConfig{
		Host:         c.cfg.Authentication.Host,
		Realm:        c.cfg.Authentication.Realm,
		ClientID:     c.cfg.Authentication.ClientID,
		ClientSecret: c.cfg.Authentication.ClientSecret,
		WebClientID:  c.cfg.Authentication.WebClientID,
	})
	return
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

func (c *Core) Start(ctx context.Context) error {
	return c.server.Start(ctx)
}
