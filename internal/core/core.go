package core

import (
	"backend/pkg/config"
	"backend/shared/logger"
	"backend/shared/logger/openob"
	"fmt"

	"context"

	"backend/internal/interfaces/http"
	kcCore "backend/shared/keycloak"

	// wshandler "backend/internal/interfaces/http/handlers/websocket"
	// appToken "backend/shared/token/app"
	// appWS "backend/shared/websocket/app"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Core struct {
	cfg    *config.Config
	logger logger.Logger
	db     *gorm.DB
	auth   *kcCore.Keycloak
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
		Headers:       c.cfg.Telemetry.User + ":" + c.cfg.Telemetry.Password,
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
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.cfg.Database.Host,
		c.cfg.Database.Port,
		c.cfg.Database.Username,
		c.cfg.Database.Password,
		c.cfg.Database.Database,
		c.cfg.Database.SSLMode)
	dialector := postgres.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(c.cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(c.cfg.Database.ConnMaxLifetime)

	// Auto migrate if enabled
	if c.cfg.Database.AutoMigrate {
		// Perform migrations here
		// Example: err = db.AutoMigrate(&models.YourModel{})
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

func (c *Core) Start(ctx context.Context) error {
	return c.server.Start(ctx)
}
