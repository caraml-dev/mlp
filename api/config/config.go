package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gojek/mlp/api/models/v2"
	"github.com/iancoleman/strcase"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	APIHost       string `validate:"required"`
	EncryptionKey string `validate:"required"`
	Environment   string `validate:"required"`
	Port          int    `validate:"required"`
	SentryDSN     string
	OauthClientID string

	Streams Streams `validate:"dive,required"`
	Docs    Documentations

	Applications  []models.Application `validate:"dive"`
	Authorization *AuthorizationConfig `validate:"required"`
	Database      *DatabaseConfig      `validate:"required"`
	Mlflow        *MlflowConfig        `validate:"required"`
	UI            *UIConfig
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

type Streams map[string][]string

type DatabaseConfig struct {
	Host          string `validate:"required"`
	Port          int    `validate:"required"`
	User          string `validate:"required"`
	Password      string `validate:"required"`
	Database      string `validate:"required"`
	MigrationPath string `validate:"required,url"`
}

type AuthorizationConfig struct {
	Enabled       bool
	KetoServerURL string `validate:"required_if=Enabled True"`
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
	FeastUIHomepage     string `json:"REACT_APP_FEAST_UI_HOMEPAGE"`
	KubeflowUIHomepage  string `json:"REACT_APP_KUBEFLOW_UI_HOMEPAGE"`
	MerlinUIHomepage    string `json:"REACT_APP_MERLIN_UI_HOMEPAGE"`
	TuringUIHomepage    string `json:"REACT_APP_TURING_UI_HOMEPAGE"`
}

// Transform env variables to the format consumed by koanf.
// The variable key is split by the double underscore ('__') sequence,
// which separates nested config variables, and then each config key is
// converted to camel-case.
//
// Example:
//
//	MY_VARIABLE => MyVariable
//	VARIABLES__ANOTHER_VARIABLE => Variables.AnotherVariable
func envVarKeyTransformer(s string) string {
	parts := strings.Split(strings.ToLower(s), "__")
	transformed := make([]string, len(parts))
	for idx, key := range parts {
		transformed[idx] = strcase.ToCamel(key)
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

	Streams: Streams{},
	Docs:    Documentations{},

	Authorization: &AuthorizationConfig{
		Enabled: false,
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
		IndexPath:  "index.html",
		StaticPath: "ui/build",
	},
}
