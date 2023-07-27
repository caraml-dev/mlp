package config_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/models"
	modelsv2 "github.com/caraml-dev/mlp/api/models/v2"
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
	oneSecond, _ := time.ParseDuration("1s")
	twoSeconds, _ := time.ParseDuration("2s")

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
					Caching: &config.InMemoryCacheConfig{
						KeyExpirySeconds:            600,
						CacheCleanUpIntervalSeconds: 900,
					},
				},
				Database: &config.DatabaseConfig{
					Host:            "localhost",
					Port:            5432,
					Database:        "mlp",
					User:            "mlp",
					MigrationPath:   "file://db-migrations",
					ConnMaxIdleTime: oneSecond,
					ConnMaxLifetime: twoSeconds,
					MaxIdleConns:    10,
					MaxOpenConns:    20,
				},
				Mlflow:       &config.MlflowConfig{},
				Docs:         []config.Documentation{},
				Applications: []modelsv2.Application{},
				Streams: map[string][]string{
					"stream-1":     {"team-a", "team-b"},
					"SecondStream": {"MyTeam"},
					"EmptyStream":  {},
				},
				UI: &config.UIConfig{
					StaticPath: "ui/build",
					IndexPath:  "index.html",
				},
				DefaultSecretStorage: &config.SecretStorage{
					Name: "default-secret-storage",
					Type: "vault",
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							Role:        "my-role",
							MountPath:   "secret",
							PathPrefix:  "caraml-secret/{{ .project }}/",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
						},
					},
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
				APIHost:     "http://localhost:8080",
				Port:        8080,
				Environment: "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled: false,
					Caching: &config.InMemoryCacheConfig{
						KeyExpirySeconds:            600,
						CacheCleanUpIntervalSeconds: 900,
					},
				},
				Database: &config.DatabaseConfig{
					Host:            "localhost",
					Port:            5432,
					Database:        "mlp",
					User:            "mlp",
					Password:        "secret",
					MigrationPath:   "file://db-migrations",
					ConnMaxIdleTime: oneSecond,
					ConnMaxLifetime: twoSeconds,
					MaxIdleConns:    10,
					MaxOpenConns:    20,
				},
				Mlflow:       &config.MlflowConfig{},
				Docs:         []config.Documentation{},
				Applications: []modelsv2.Application{},
				Streams: map[string][]string{
					"stream-1":     {"team-a", "team-b"},
					"SecondStream": {"MyTeam"},
					"EmptyStream":  {},
				},
				UI: &config.UIConfig{
					StaticPath: "ui/build",
					IndexPath:  "index.html",
				},
				DefaultSecretStorage: &config.SecretStorage{
					Name: "default-secret-storage",
					Type: "vault",
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							Role:        "my-role",
							MountPath:   "secret",
							PathPrefix:  "caraml-secret/{{ .project }}/",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
						},
					},
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
				Environment:   "dev",
				OauthClientID: "oauth-client-id",
				SentryDSN:     "1234",
				Applications: []modelsv2.Application{
					{
						Name:        "Turing",
						Description: "ML Experimentation System",
						Homepage:    "/turing",
						Configuration: &modelsv2.ApplicationConfig{
							API:      "/api/turing/v1",
							IconName: "graphApp",
							Navigation: []modelsv2.NavigationMenuItem{
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
					Enabled:         true,
					KetoRemoteRead:  "http://localhost:4466",
					KetoRemoteWrite: "http://localhost:4467",
					Caching: &config.InMemoryCacheConfig{
						Enabled:                     true,
						KeyExpirySeconds:            1000,
						CacheCleanUpIntervalSeconds: 2000,
					},
					UseMiddleware: true,
				},
				Database: &config.DatabaseConfig{
					Host:            "localhost",
					Port:            5432,
					Database:        "mlp",
					User:            "mlp",
					Password:        "secret",
					MigrationPath:   "file://db-migrations",
					ConnMaxIdleTime: oneSecond,
					ConnMaxLifetime: twoSeconds,
					MaxIdleConns:    10,
					MaxOpenConns:    20,
				},
				Docs: []config.Documentation{
					{
						Label: "Merlin User Guide",
						Href:  "https://github.com/caraml-dev/merlin/blob/main/docs/getting-started/README.md",
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
				DefaultSecretStorage: &config.SecretStorage{
					Name: "default-secret-storage",
					Type: "vault",
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							Role:        "my-role",
							MountPath:   "secret",
							PathPrefix:  "caraml-secret/{{ .project }}/",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
						},
					},
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
				APIHost:     "/v1",
				Port:        8080,
				Environment: "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled: false,
					Caching: &config.InMemoryCacheConfig{
						KeyExpirySeconds:            600,
						CacheCleanUpIntervalSeconds: 900,
					},
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
				DefaultSecretStorage: &config.SecretStorage{
					Name: "default-secret-storage",
					Type: "vault",
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							Role:        "my-role",
							MountPath:   "secret",
							PathPrefix:  "caraml-secret/{{ .project }}/",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
						},
					},
				},
			},
		},
		"extended | success": {
			config: &config.Config{
				APIHost:     "/v1",
				Port:        8080,
				Environment: "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled:         true,
					KetoRemoteRead:  "http://keto.mlp",
					KetoRemoteWrite: "http://keto.mlp",
					Caching: &config.InMemoryCacheConfig{
						KeyExpirySeconds:            600,
						CacheCleanUpIntervalSeconds: 900,
					},
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
				DefaultSecretStorage: &config.SecretStorage{
					Name: "default-secret-storage",
					Type: "vault",
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							Role:        "my-role",
							MountPath:   "secret",
							PathPrefix:  "caraml-secret/{{ .project }}/",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
						},
					},
				},
			},
		},
		"default config | failure": {
			config: config.NewDefaultConfig(),
			error: errors.New(
				"failed to validate configuration: " +
					"Key: 'Config.Database.User' Error:Field validation for 'User' failed on the 'required' tag\n" +
					"Key: 'Config.Database.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			),
		},
		"missing authz server | failure": {
			config: &config.Config{
				APIHost:     "/v1",
				Port:        8080,
				Environment: "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled:         true,
					KetoRemoteWrite: "localhost:4467",
					Caching: &config.InMemoryCacheConfig{
						KeyExpirySeconds:            600,
						CacheCleanUpIntervalSeconds: 900,
					},
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
				DefaultSecretStorage: &config.SecretStorage{
					Name: "default-secret-storage",
					Type: "vault",
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							Role:        "my-role",
							MountPath:   "secret",
							PathPrefix:  "caraml-secret/{{ .project }}/",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
						},
					},
				},
			},
			error: errors.New(
				"failed to validate configuration: " +
					"Key: 'Config.Authorization.KetoRemoteRead' " +
					"Error:Field validation for 'KetoRemoteRead' failed on the 'required_if' tag",
			),
		},
		"missing authz cache key expiry | failure": {
			config: &config.Config{
				APIHost:     "/v1",
				Port:        8080,
				Environment: "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled:         true,
					KetoRemoteRead:  "http://abc",
					KetoRemoteWrite: "http://abc",
					Caching: &config.InMemoryCacheConfig{
						Enabled: true,
					},
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
				DefaultSecretStorage: &config.SecretStorage{
					Name: "default-secret-storage",
					Type: "vault",
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							Role:        "my-role",
							MountPath:   "secret",
							PathPrefix:  "caraml-secret/{{ .project }}/",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
						},
					},
				},
			},
			error: errors.New(
				"failed to validate configuration: " +
					"Key: 'Config.Authorization.Caching.KeyExpirySeconds' " +
					"Error:Field validation for 'KeyExpirySeconds' failed on the 'required_if' tag\n" +
					"Key: 'Config.Authorization.Caching.CacheCleanUpIntervalSeconds' " +
					"Error:Field validation for 'CacheCleanUpIntervalSeconds' failed on the 'required_if' tag",
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
