package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type TenantConfiguration struct {
	WebClient     *Client        `yaml:"web-client" mapstructure:"web-client"`
	BackendClient *Client        `yaml:"backend-client" mapstructure:"backend-client"`
	Authorization *Authorization `yaml:"authorization" mapstructure:"authorization"`
}

type Client struct {
	Type           string   `yaml:"type" mapstructure:"type"`
	ClientID       string   `yaml:"client_id" mapstructure:"client_id"`
	Name           string   `yaml:"name" mapstructure:"name"`
	Description    string   `yaml:"description" mapstructure:"description"`
	RedirectURIs   []string `yaml:"redirect_uris" mapstructure:"redirect_uris"`
	Origins        []string `yaml:"origins" mapstructure:"origins"`
	LogoutRedirect *string  `yaml:"logout_redirect,omitempty" mapstructure:"logout_redirect,omitempty"`
	RootURL        *string  `yaml:"root_url,omitempty" mapstructure:"root_url,omitempty"`
	AdminURL       *string  `yaml:"admin_url,omitempty" mapstructure:"admin_url,omitempty"`
}

type Authorization struct {
	AdminGroup string              `yaml:"roles" mapstructure:"admin_group"`
	Roles      []string            `yaml:"roles" mapstructure:"roles"`
	Resources  map[string]Resource `yaml:"resources" mapstructure:"resources"`
}

type Resource struct {
	Name        string              `yaml:"name" mapstructure:"name"`
	Type        string              `yaml:"type" mapstructure:"type"`
	Scopes      *[]string           `yaml:"type" mapstructure:"scopes"` //Resources scopes allowed if empty all
	Permissions map[string][]string `yaml:"permissions,omitempty" mapstructure:"permissions,omitempty"`
}

func LoadTenant() (*TenantConfiguration, error) {
	var config TenantConfiguration

	v := viper.New()

	fileName := os.Getenv("AUTH_CONFIG_FILE_NAME")
	fileType := os.Getenv("AUTH_CONFIG_FILE_TYPE")
	path := os.Getenv("AUTH_CONFIG_FILE_PATH")

	if fileName == "" || fileType == "" || path == "" {
		return nil, fmt.Errorf("missing required environment variables: AUTH_CONFIG_FILE_NAME, AUTH_CONFIG_FILE_TYPE, or AUTH_CONFIG_FILE_PATH")
	}

	v.SetConfigName(fileName)
	v.SetConfigType(fileType)
	v.AddConfigPath(path)

	v.SetEnvPrefix("AUTH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading auth config file: %w", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling auth config: %w", err)
	}

	return &config, nil
}
