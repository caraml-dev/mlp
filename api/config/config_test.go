package config_test

import (
	"errors"
	"os"
	"testing"

	"github.com/gojek/mlp/api/config"
	"github.com/gojek/mlp/api/models/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func envSetter(envs map[string]string) (closer func()) {
	originalEnvs := map[string]string{}

	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvs[name] = originalValue
		}
		_ = os.Setenv(name, value)
	}

	return func() {
		for name := range envs {
			origValue, has := originalEnvs[name]
			if has {
				_ = os.Setenv(name, origValue)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}
}

func TestLoad(t *testing.T) {
	suite := map[string]struct {
		configs  []string
		env      map[string]string
		expected *config.Config
		error    error
	}{
		"default | success": {
			configs:  []string{},
			env:      map[string]string{},
			expected: config.NewDefaultConfig(),
		},
		"config-1.yaml | success": {
			configs: []string{"testdata/config-1.yaml"},
			env:     map[string]string{},
			expected: &config.Config{
				APIHost:     "http://localhost:8080",
				Port:        8080,
				Environment: "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled: false,
				},
				Database: &config.DatabaseConfig{
					Host:          "localhost",
					Port:          5432,
					Database:      "mlp",
					User:          "mlp",
					MigrationPath: "file://db-migrations",
				},
				Mlflow: &config.MlflowConfig{},
				Docs:   []config.Documentation{},
				Streams: map[string][]string{
					"stream-1":     {"team-a", "team-b"},
					"SecondStream": {"MyTeam"},
					"EmptyStream":  {},
				},
				UI: &config.UIConfig{
					StaticPath: "ui/build",
					IndexPath:  "index.html",
				},
			},
		},
		"config-1.yaml + env variables | success": {
			configs: []string{"testdata/config-1.yaml"},
			env: map[string]string{
				"ENCRYPTION_KEY":     "test-key",
				"DATABASE__PASSWORD": "secret",
			},
			expected: &config.Config{
				APIHost:       "http://localhost:8080",
				Port:          8080,
				EncryptionKey: "test-key",
				Environment:   "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled: false,
				},
				Database: &config.DatabaseConfig{
					Host:          "localhost",
					Port:          5432,
					Database:      "mlp",
					User:          "mlp",
					Password:      "secret",
					MigrationPath: "file://db-migrations",
				},
				Mlflow: &config.MlflowConfig{},
				Docs:   []config.Documentation{},
				Streams: map[string][]string{
					"stream-1":     {"team-a", "team-b"},
					"SecondStream": {"MyTeam"},
					"EmptyStream":  {},
				},
				UI: &config.UIConfig{
					StaticPath: "ui/build",
					IndexPath:  "index.html",
				},
			},
		},
		"config-1.yaml + config-2.yaml + env variables | success": {
			configs: []string{"testdata/config-1.yaml", "testdata/config-2.yaml"},
			env: map[string]string{
				"ENCRYPTION_KEY":       "test-key",
				"OAUTH_CLIENT_ID":      "oauth-client-id",
				"DATABASE__PASSWORD":   "secret",
				"MLFLOW__TRACKING_URL": "http://mlflow.dev",
				"SENTRY_DSN":           "1234",
			},
			expected: &config.Config{
				APIHost:       "http://localhost:8080",
				Port:          8080,
				EncryptionKey: "test-key",
				Environment:   "dev",
				OauthClientID: "oauth-client-id",
				SentryDSN:     "1234",
				Applications: []models.Application{
					{
						Name:        "Turing",
						Description: "ML Experimentation System",
						Homepage:    "/turing",
						Configuration: &models.ApplicationConfig{
							API:      "/api/turing/v1",
							IconName: "graphApp",
							Navigation: []models.NavigationMenuItem{
								{
									Label:       "Routers",
									Destination: "/routers",
								},
								{
									Label:       "Experiments",
									Destination: "/experiments",
								},
							},
						},
					},
				},
				Authorization: &config.AuthorizationConfig{
					Enabled:       true,
					KetoServerURL: "http://localhost:4466",
				},
				Database: &config.DatabaseConfig{
					Host:          "localhost",
					Port:          5432,
					Database:      "mlp",
					User:          "mlp",
					Password:      "secret",
					MigrationPath: "file://db-migrations",
				},
				Docs: []config.Documentation{
					{
						Label: "Merlin User Guide",
						Href:  "https://github.com/gojek/merlin/blob/main/docs/getting-started/README.md",
					},
				},
				Mlflow: &config.MlflowConfig{
					TrackingURL: "http://mlflow.dev",
				},
				Streams: map[string][]string{
					"stream-1":     {"team-a", "team-b"},
					"SecondStream": {"MyTeam"},
					"EmptyStream":  {},
				},
				UI: &config.UIConfig{
					StaticPath: "ui/build",
					IndexPath:  "index.html",

					ClockworkUIHomepage: "http://clockwork.dev",
					KubeflowUIHomepage:  "http://kubeflow.org",
				},
			},
		},
		"config files doesn't exist | failure": {
			configs: []string{"invalid-file"},
			error: errors.New(
				"failed to read config from file 'invalid-file': open invalid-file: no such file or directory",
			),
		},
	}

	for name, tt := range suite {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(envSetter(tt.env))
			actual, err := config.Load(tt.configs...)

			if tt.error == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			} else {
				assert.EqualError(t, err, tt.error.Error())
			}
		})
	}
}

func TestValidate(t *testing.T) {
	suite := map[string]struct {
		config *config.Config
		error  error
	}{
		"minimal | success": {
			config: &config.Config{
				APIHost:       "/v1",
				Port:          8080,
				Environment:   "dev",
				EncryptionKey: "secret-key",
				Authorization: &config.AuthorizationConfig{
					Enabled: false,
				},
				Database: &config.DatabaseConfig{
					Host:          "localhost",
					Port:          5432,
					User:          "mlp",
					Password:      "mlp",
					Database:      "mlp",
					MigrationPath: "file://db-migrations",
				},
				Mlflow: &config.MlflowConfig{
					TrackingURL: "http://mlflow.tracking",
				},
			},
		},
		"extended | success": {
			config: &config.Config{
				APIHost:       "/v1",
				Port:          8080,
				Environment:   "dev",
				EncryptionKey: "secret-key",
				Authorization: &config.AuthorizationConfig{
					Enabled:       true,
					KetoServerURL: "http://keto.mlp",
				},
				Database: &config.DatabaseConfig{
					Host:          "localhost",
					Port:          5432,
					User:          "mlp",
					Password:      "mlp",
					Database:      "mlp",
					MigrationPath: "file://db-migrations",
				},
				Mlflow: &config.MlflowConfig{
					TrackingURL: "http://mlflow.tracking",
				},
				Streams: map[string][]string{
					"my-stream": {"my-team"},
				},
			},
		},
		"default config | failure": {
			config: config.NewDefaultConfig(),
			error: errors.New(
				"failed to validate configuration: " +
					"Key: 'Config.EncryptionKey' Error:Field validation for 'EncryptionKey' failed on the 'required' tag\n" +
					"Key: 'Config.Database.User' Error:Field validation for 'User' failed on the 'required' tag\n" +
					"Key: 'Config.Database.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			),
		},
		"missing auth server | failure": {
			config: &config.Config{
				APIHost:       "/v1",
				Port:          8080,
				Environment:   "dev",
				EncryptionKey: "secret-key",
				Authorization: &config.AuthorizationConfig{
					Enabled: true,
				},
				Database: &config.DatabaseConfig{
					Host:          "localhost",
					Port:          5432,
					User:          "mlp",
					Password:      "mlp",
					Database:      "mlp",
					MigrationPath: "file://db-migrations",
				},
				Mlflow: &config.MlflowConfig{
					TrackingURL: "http://mlflow.tracking",
				},
			},
			error: errors.New(
				"failed to validate configuration: " +
					"Key: 'Config.Authorization.KetoServerURL' " +
					"Error:Field validation for 'KetoServerURL' failed on the 'required_if' tag",
			),
		},
	}

	for name, tt := range suite {
		t.Run(name, func(t *testing.T) {
			actual := config.Validate(tt.config)

			if tt.error == nil {
				require.NoError(t, actual)
			} else {
				assert.EqualError(t, actual, tt.error.Error())
			}
		})
	}
}
