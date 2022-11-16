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

func Load(paths ...string) (*Config, error) {
	k := koanf.New(".")

	for _, f := range paths {
		err := k.Load(file.Provider(f), yaml.Parser())
		if err != nil {
			return nil, fmt.Errorf("failed to read config from file '%s': %s", f, err)
		}
	}

	err := k.Load(env.Provider("", ".", func(s string) string {
		parts := strings.Split(strings.ToLower(s), "::")
		transformed := make([]string, len(parts))
		for idx, key := range parts {
			transformed[idx] = strcase.ToCamel(key)
		}

		return strings.Join(transformed, ".")
	}), nil)

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
