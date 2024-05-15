package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"

	"github.com/caraml-dev/mlp/api/models"
	modelsv2 "github.com/caraml-dev/mlp/api/models/v2"
)

type Config struct {
	APIHost       string `validate:"required"`
	Environment   string `validate:"required"`
	Port          int    `validate:"required"`
	SentryDSN     string
	OauthClientID string

	Streams Streams `validate:"dive,required"`
	Docs    Documentations

	Applications         []modelsv2.Application `validate:"dive"`
	Authorization        *AuthorizationConfig   `validate:"required"`
	Database             *DatabaseConfig        `validate:"required"`
	Mlflow               *MlflowConfig          `validate:"required"`
	DefaultSecretStorage *SecretStorage         `validate:"required"`
	UI                   *UIConfig
}

// SecretStorage represents the configuration for a secret storage.
type SecretStorage struct {
	// Name is the name of the secret storage.
	Name string `validate:"required"`
	// Type is the type of the secret storage.
	Type string `validate:"oneof=internal vault"`
	// Config is the configuration of the secret storage.
	Config models.SecretStorageConfig
}

func NewDefaultConfig() *Config {
	var c Config

	// make a deep copy of the default config via JSON serialization
	defaultBytes, _ := json.Marshal(defaultConfig)
	_ = json.Unmarshal(defaultBytes, &c)

	return &c
}

func (c *Config) ListenAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}

// DefaultSecretStorageModel returns the default secret storage model from the given config.
// The returned secret storage model is a globally-scoped secret storage.
func (c *Config) DefaultSecretStorageModel() *models.SecretStorage {
	return &models.SecretStorage{
		Name:   c.DefaultSecretStorage.Name,
		Type:   models.SecretStorageType(c.DefaultSecretStorage.Type),
		Scope:  models.GlobalSecretStorageScope,
		Config: c.DefaultSecretStorage.Config,
	}
}

type Streams map[string][]string

type DatabaseConfig struct {
	Host          string `validate:"required"`
	Port          int    `validate:"required"`
	User          string `validate:"required"`
	Password      string `validate:"required"`
	Database      string `validate:"required"`
	MigrationPath string `validate:"required,url"`

	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

type AuthorizationConfig struct {
	Enabled         bool
	KetoRemoteRead  string               `validate:"required_if=Enabled True"`
	KetoRemoteWrite string               `validate:"required_if=Enabled True"`
	Caching         *InMemoryCacheConfig `validate:"required_if=Enabled True"`
	UseMiddleware   bool
}

type InMemoryCacheConfig struct {
	Enabled                     bool
	KeyExpirySeconds            int `validate:"required_if=Enabled True"`
	CacheCleanUpIntervalSeconds int `validate:"required_if=Enabled True"`
}

type MlflowConfig struct {
	TrackingURL string `validated:"required,url"`
}

type Documentations []Documentation

type Documentation struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

// UIConfig stores the configuration for the UI.
type UIConfig struct {
	StaticPath string `validated:"required"`
	IndexPath  string `validated:"required"`

	ClockworkUIHomepage string `json:"REACT_APP_CLOCKWORK_UI_HOMEPAGE"`
	KubeflowUIHomepage  string `json:"REACT_APP_KUBEFLOW_UI_HOMEPAGE"`

	AllowCustomStream bool `json:"REACT_APP_ALLOW_CUSTOM_STREAM"`
	AllowCustomTeam   bool `json:"REACT_APP_ALLOW_CUSTOM_TEAM"`
}

// Transform env variables to the format consumed by koanf.
// The variable key is split by the double underscore ('__') sequence,
// which separates nested config variables, and then each config key is
// converted to lower camel-case.
//
// Example:
//
//	MY_VARIABLE => MyVariable
//	VARIABLES__ANOTHER_VARIABLE => variables.anotherVariable
func envVarKeyTransformer(s string) string {
	parts := strings.Split(strings.ToLower(s), "__")
	transformed := make([]string, len(parts))
	for idx, key := range parts {
		transformed[idx] = strcase.ToLowerCamel(key)
	}

	return strings.Join(transformed, ".")
}

func Load(paths ...string) (*Config, error) {
	k := koanf.New(".")

	// read config from zero or more YAML config files
	for _, f := range paths {
		err := k.Load(file.Provider(f), yaml.Parser())
		if err != nil {
			return nil, fmt.Errorf("failed to read config from file '%s': %s", f, err)
		}
	}

	// read config overrides from env variables
	err := k.Load(env.Provider("", ".", envVarKeyTransformer), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from environment variables: %s", err)
	}

	// create config instance with pre-populated default values
	config := NewDefaultConfig()

	err = k.Unmarshal("", config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall config values: %s", err)
	}

	return config, err
}

func Validate(config *Config) error {
	validate := validator.New()

	err := validate.Struct(config)
	if err != nil {
		return fmt.Errorf("failed to validate configuration: %s", err)
	}
	return nil
}

func LoadAndValidate(paths ...string) (*Config, error) {
	config, err := Load(paths...)
	if err != nil {
		return nil, err
	}

	err = Validate(config)
	return config, err
}

var defaultConfig = &Config{
	APIHost:     "http://localhost:8080",
	Environment: "dev",
	Port:        8080,

	Streams:      Streams{},
	Docs:         Documentations{},
	Applications: []modelsv2.Application{},
	Authorization: &AuthorizationConfig{
		Enabled: false,
		Caching: &InMemoryCacheConfig{
			KeyExpirySeconds:            600,
			CacheCleanUpIntervalSeconds: 900,
		},
		UseMiddleware: false,
	},
	Database: &DatabaseConfig{
		Host:          "localhost",
		Port:          5432,
		Database:      "mlp",
		MigrationPath: "file://db-migrations",
	},
	Mlflow: &MlflowConfig{
		TrackingURL: "",
	},
	UI: &UIConfig{
		IndexPath:         "index.html",
		StaticPath:        "ui/build",
		AllowCustomTeam:   true,
		AllowCustomStream: true,
	},
	DefaultSecretStorage: &SecretStorage{
		Name: "internal",
		Type: "internal",
	},
}
