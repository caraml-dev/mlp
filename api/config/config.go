package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

type Config struct {
	APIHost       string `envconfig:"API_HOST" default:"http://localhost:8080/v1`
	Port          int    `envconfig:"PORT" default:"8080"`
	Environment   string `envconfig:"ENVIRONMENT" default:"dev"`
	EncryptionKey string `envconfig:"ENCRYPTION_KEY" required:"true"`

	MlflowConfig        MlflowConfig
	DbConfig            DatabaseConfig
	GitlabConfig        GitlabConfig
	AuthorizationConfig AuthorizationConfig
	UI                  UIConfig

	OauthClientID string `envconfig:"OAUTH_CLIENT_ID"`
	SentryDSN     string `envconfig:"SENTRY_DSN"`
}

// UIConfig stores the configuration for the UI.
type UIConfig struct {
	StaticPath string `envconfig:"UI_STATIC_PATH" default:"ui/build"`
	IndexPath  string `envconfig:"UI_INDEX_PATH" default:"index.html"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DATABASE_HOST" required:"true"`
	Port     int    `envconfig:"DATABASE_PORT" default:"5432"`
	User     string `envconfig:"DATABASE_USER" required:"true"`
	Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
	Database string `envconfig:"DATABASE_NAME" default:"mlp"`
}

type GitlabConfig struct {
	Enabled      bool     `envconfig:"GITLAB_ENABLED" default:"false"`
	Host         string   `envconfig:"GITLAB_HOST"`
	ClientID     string   `envconfig:"GITLAB_CLIENT_ID"`
	ClientSecret string   `envconfig:"GITLAB_CLIENT_SECRET"`
	RedirectURL  string   `envconfig:"GITLAB_REDIRECT_URL"`
	Scopes       []string `envconfig:"GITLAB_OAUTH_SCOPES" default:"read_user"`
}

type AuthorizationConfig struct {
	AuthorizationEnabled   bool   `envconfig:"AUTHORIZATION_ENABLED" default:"false"`
	AuthorizationServerUrl string `envconfig:"AUTHORIZATION_SERVER_URL" default:"http://localhost:4466"`
}

type MlflowConfig struct {
	TrackingUrl string `envconfig:"MLFLOW_TRACKING_URL" required:"true"`
}

func InitConfigEnv() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *GitlabConfig) InitOauthConfig() (*oauth2.Config, error) {
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       cfg.Scopes,
		RedirectURL:  cfg.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", cfg.Host),
			TokenURL: fmt.Sprintf("%s/oauth/token", cfg.Host),
		},
	}, nil
}

func (c *Config) ListenAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}
