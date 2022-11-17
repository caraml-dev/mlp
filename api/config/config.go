package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

const EnvVarPrefix = "CARAML_"

type Config struct {
	APIHost       string
	EncryptionKey string
	Environment   string
	Port          int
	SentryDSN     string
	OauthClientID string

	Streams Streams
	Docs    Documentations

	Authorization *AuthorizationConfig
	Database      *DatabaseConfig
	Mlflow        *MlflowConfig
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
	Host          string
	Port          int
	User          string
	Password      string
	Database      string
	MigrationPath string
}

type AuthorizationConfig struct {
	Enabled       bool
	KetoServerURL string
}

type MlflowConfig struct {
	TrackingURL string
}

type Documentations []Documentation

type Documentation struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

// UIConfig stores the configuration for the UI.
type UIConfig struct {
	StaticPath string
	IndexPath  string

	FeastCoreAPI        string `json:"REACT_APP_FEAST_CORE_API"`
	MerlinAPI           string `json:"REACT_APP_MERLIN_API"`
	TuringAPI           string `json:"REACT_APP_TURING_API"`
	ClockworkUIHomepage string `json:"REACT_APP_CLOCKWORK_UI_HOMEPAGE"`
	FeastUIHomepage     string `json:"REACT_APP_FEAST_UI_HOMEPAGE"`
	KubeflowUIHomepage  string `json:"REACT_APP_KUBEFLOW_UI_HOMEPAGE"`
	MerlinUIHomepage    string `json:"REACT_APP_MERLIN_UI_HOMEPAGE"`
	TuringUIHomepage    string `json:"REACT_APP_TURING_UI_HOMEPAGE"`
}

// Transform env variables to the format consumed by koanf.
// First, "CARAML_" prefix is trimmed from the variable key, then it's
// split by the double underscore ('__') sequence, which separates nested
// config variables, and then each config key is converted to camel-case.
//
// Example:
//	CARAML_MY_VARIABLE => MyVariable
//	CARAML_VARIABLES__ANOTHER_VARIABLE => Variables.AnotherVariable
func envVarKeyTransformer(s string) string {
	parts := strings.Split(strings.ToLower(strings.TrimPrefix(s, EnvVarPrefix)), "__")
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
	err := k.Load(env.Provider(EnvVarPrefix, ".", envVarKeyTransformer), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from environment variables: %s", err)
	}

	// create config instance with pre-populated default values
	config := NewDefaultConfig()
	err = k.Unmarshal("", config)

	return config, err
}

var defaultConfig = &Config{
	APIHost:     "http://localhost:8080/v1",
	Environment: "dev",
	Port:        8080,

	Streams: Streams{},
	Docs:    Documentations{},

	Authorization: &AuthorizationConfig{
		Enabled:       false,
		KetoServerURL: "http://localhost:4466",
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
