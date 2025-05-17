package config

import (
	"fmt"
	"time"
)

type Config struct {
	App            AppConfig            `mapstructure:"app"`
	HTTP           HTTPConfig           `mapstructure:"http"`
	GRPC           GRPCConfig           `mapstructure:"grpc"`
	Database       DatabaseConfig       `mapstructure:"database"`
	Stream         StreamConfig         `mapstructure:"stream"`
	Websocket      WebsocketConfig      `mapstructure:"websocket"`
	Telemetry      TelemetryConfig      `mapstructure:"telemetry"`
	Authentication AuthenticationConfig `mapstructure:"authentication"`
	Ingest         IngestConfig         `mapstructure:"ingest"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"log_level"`
}

type HTTPConfig struct {
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	ShutdownTimeout    time.Duration `mapstructure:"shutdown_timeout"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
	CORSAllowedOrigins []string      `mapstructure:"cors_allowed_origins"`
}

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	AutoMigrate     bool          `mapstructure:"auto_migrate"`
}

type IngestConfig struct {
	URL          string `mapstructure:"url"`
	ClientID     string `mapstructure:"client_id"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	ManageURL    string `mapstructure:"manage_url"`
	ManageKey    string `mapstructure:"manage_key"`
	ManageSecret string `mapstructure:"manage_secret"`
}
type StreamConfig struct {
	URL     string        `mapstructure:"url"`
	Cluster string        `mapstructure:"cluster"`
	Client  string        `mapstructure:"client"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type WebsocketConfig struct {
	ReadBufferSize   int           `mapstructure:"read_buffer_size"`
	WriteBufferSize  int           `mapstructure:"write_buffer_size"`
	HandshakeTimeout time.Duration `mapstructure:"handshake_timeout"`
	PingInterval     time.Duration `mapstructure:"ping_interval"`
	MaxMessageSize   int64         `mapstructure:"max_message_size"`
}

type TelemetryConfig struct {
	ServiceName    string  `mapstructure:"service_name"`
	Environment    string  `mapstructure:"enviroment"`
	Endpoint       string  `mapstructure:"endpoint"`
	TracingEnabled bool    `mapstructure:"tracing_enabled"`
	SamplingRate   float64 `mapstructure:"sampling_rate"`
	User           string  `mapstructure:"user"`
	Password       string  `mapstructure:"password"`
	TenantID       string  `mapstructure:"tenant_id"`
	StreamName     string  `mapstructure:"stream_name"`
	MetricsEnabled bool    `mapstructure:"metrics_enabled"`
	MetricsHost    string  `mapstructure:"metrics_host"`
	MetricsPort    int     `mapstructure:"metrics_port"`
}

type AuthenticationConfig struct {
	Host         string `mapstructure:"host"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	WebClientID  string `mapstructure:"client_web"`
}

func (c *DatabaseConfig) DSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
	case "sqlite":
		return c.Database
	default:
		return ""
	}
}
