package config_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gojek/mlp/api/config"
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
				APIHost:     "http://localhost:8080/v1",
				Port:        8080,
				Environment: "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled:       false,
					KetoServerURL: "http://localhost:4466",
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
				"CARAML_ENCRYPTION_KEY":     "test-key",
				"CARAML_DATABASE__PASSWORD": "secret",
			},
			expected: &config.Config{
				APIHost:       "http://localhost:8080/v1",
				Port:          8080,
				EncryptionKey: "test-key",
				Environment:   "dev",
				Authorization: &config.AuthorizationConfig{
					Enabled:       false,
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
				"CARAML_ENCRYPTION_KEY":       "test-key",
				"CARAML_OAUTH_CLIENT_ID":      "oauth-client-id",
				"CARAML_DATABASE__PASSWORD":   "secret",
				"CARAML_MLFLOW__TRACKING_URL": "http://mlflow.dev",
				"CARAML_SENTRY_DSN":           "1234",
			},
			expected: &config.Config{
				APIHost:       "http://localhost:8080/v1",
				Port:          8080,
				EncryptionKey: "test-key",
				Environment:   "dev",
				OauthClientID: "oauth-client-id",
				SentryDSN:     "1234",
				Authorization: &config.AuthorizationConfig{
					Enabled:       false,
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

					FeastCoreAPI:        "/feast/api",
					MerlinAPI:           "/api/merlin/v1",
					TuringAPI:           "/api/turing/v1",
					ClockworkUIHomepage: "http://clockwork.dev",
					FeastUIHomepage:     "/feast",
					KubeflowUIHomepage:  "http://kubeflow.org",
					MerlinUIHomepage:    "/merlin",
					TuringUIHomepage:    "/turing",
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

			fmt.Printf("############ %s:\n\t%s\n", name, strings.Join(os.Environ(), "\n\t"))

			if tt.error == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			} else {
				assert.EqualError(t, err, tt.error.Error())
			}
		})
	}
}
