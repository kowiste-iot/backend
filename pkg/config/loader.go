package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Load loads the configuration from various sources
func Load() (*Config, error) {
	var config Config

	v := viper.New()

	// Set default configurations
	setDefaults(v)
	fileName := os.Getenv("CONFIG_FILE_NAME")
	fileType := os.Getenv("CONFIG_FILE_TYPE")
	path := os.Getenv("CONFIG_FILE_PATH")
	// Configure Viper
	v.SetConfigName(fileName)
	v.SetConfigType(fileType)
	v.AddConfigPath(path)
	v.AddConfigPath(".") // optionally look for config in working directory

	// Enable environment variables
	v.SetEnvPrefix("APP") // prefix for environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "asset-service")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.log_level", "info")

	// HTTP defaults
	v.SetDefault("http.host", "0.0.0.0")
	v.SetDefault("http.port", 5000)
	v.SetDefault("http.shutdown_timeout", "5s")
	v.SetDefault("http.read_timeout", "15s")
	v.SetDefault("http.write_timeout", "15s")
	v.SetDefault("http.cors_allowed_origins", []string{"*"})

	// GRPC defaults
	v.SetDefault("grpc.host", "0.0.0.0")
	v.SetDefault("grpc.port", 9090)

	// Database defaults
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.username", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.database", "asset_service")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.auto_migrate", true)

	//Ingest
	v.SetDefault("ingets.url", "tcp://localhost:1883")
	v.SetDefault("ingest.client_id", "ingest-service")
	v.SetDefault("ingets.user", "admin")
	v.SetDefault("ingest.password", "")

	// Stream defaults
	v.SetDefault("stream.url", "nats://localhost:4222")
	v.SetDefault("stream.timeout", "10s")

	// JWT defaults
	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.expiration_time", "24h")
	v.SetDefault("jwt.refresh_duration", "72h")

	// Websocket defaults
	v.SetDefault("websocket.read_buffer_size", 1024)
	v.SetDefault("websocket.write_buffer_size", 1024)
	v.SetDefault("websocket.handshake_timeout", "10s")
	v.SetDefault("websocket.ping_interval", "30s")
	v.SetDefault("websocket.max_message_size", 512000) // 512KB

	// Telemetry defaults
	v.SetDefault("telemetry.service_name", "asset-service")
	v.SetDefault("telemetry.environment", "development")
	v.SetDefault("telemetry.endpoint", "localhost:5080")
	v.SetDefault("telemetry.tracing_enabled", true)
	v.SetDefault("telemetry.sampling_rate", 1.0)
	v.SetDefault("telemetry.tenant_id", "default")
	v.SetDefault("telemetry.stream_name", "logs")
	v.SetDefault("telemetry.user", "admin")
	v.SetDefault("telemetry.password", "password")
	v.SetDefault("telemetry.metrics_enabled", true)
	v.SetDefault("telemetry.metrics_host", "0.0.0.0")
	v.SetDefault("telemetry.metrics_port", 9100)

	//Authentication defaults
	v.SetDefault("authentication.host", "http://localhost:7080/auth")
	v.SetDefault("authentication.realm", "master")
	v.SetDefault("authentication.client_id", "master-client")
	v.SetDefault("authentication.client_secret", "master_password1!")
	v.SetDefault("authentication.client_web", "vue-client")

}
