package secretstorage

import (
	"fmt"
	"testing"

	"github.com/caraml-dev/mlp/api/models"
	mlperror "github.com/caraml-dev/mlp/api/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type VaultSecretStorageClientTestSuite struct {
	suite.Suite
	client SecretStorageClient
}

func (s *VaultSecretStorageClientTestSuite) SetupSuite() {
	secretStorage := &models.SecretStorage{
		Name:  "test-storage",
		Type:  models.VaultSecretStorageType,
		Scope: models.GlobalSecretStorageScope,
		Config: models.SecretStorageConfig{
			VaultConfig: &models.VaultConfig{
				URL:        "http://localhost:8200",
				Role:       "my-role",
				MountPath:  "secret",
				PathPrefix: "caraml/{{ .Project }}",
				AuthMethod: models.TokenAuthMethod,
				Token:      "root",
			},
		},
	}

	client, err := NewSecretStorageClient(secretStorage)
	if err != nil {
		s.FailNow("failed to create secret storage vaultClient", err)
	}
	s.client = client
}

func (s *VaultSecretStorageClientTestSuite) TestGet() {
	type args struct {
		name    string
		project string
	}

	project := "test-get"
	existingSecrets := map[string]string{
		"secret_1": "value_1",
	}

	s.initializeSecrets(existingSecrets, project)

	tests := []struct {
		name      string
		args      args
		want      string
		wantError error
	}{
		{
			name: "success: get existing secret",
			args: args{
				name:    "secret_1",
				project: project,
			},
			want: existingSecrets["secret_1"],
		},
		{
			name: "failed: get non existent secret in same secret path",
			args: args{
				name:    "secret_2",
				project: project,
			},
			wantError: mlperror.NewNotFoundErrorf("secret %s not found in project %s", "secret_2", project),
		},
		{
			name: "failed: secret not found",
			args: args{
				name:    "secret_1",
				project: "test-get-2",
			},
			wantError: mlperror.NewNotFoundErrorf("secret %s not found in project %s", "secret_1", "test-get-2"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.client.Get(tt.args.name, tt.args.project)
			if tt.wantError != nil {
				s.Assert().Error(err)
				s.Assert().Equal(tt.wantError, err)
				return
			}
			s.Assert().NoError(err)
			s.Assert().Equal(tt.want, got)
		})
	}

	s.deleteSecrets(existingSecrets, project)
}

func (s *VaultSecretStorageClientTestSuite) TestSet() {
	type args struct {
		name    string
		value   string
		project string
	}

	project := "test-set"
	existingSecrets := map[string]string{
		"secret_1": "value_1",
	}

	s.initializeSecrets(existingSecrets, project)

	tests := []struct {
		name             string
		args             args
		want             string
		wantErrorMessage string
	}{
		{
			name: "success: set new secret in same project as existing secret",
			args: args{
				name:    "new_secret",
				value:   "new_value",
				project: project,
			},
			want: "new_value",
		},
		{
			name: "success: set existing secret in same project as existing secret",
			args: args{
				name:    "secret_1",
				value:   "new_value",
				project: project,
			},
			want: "new_value",
		},
		{
			name: "success: set new secret in different project as existing secret",
			args: args{
				name:    "new_secret",
				value:   "new_value",
				project: "test-set-2",
			},
			want: "new_value",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.client.Set(tt.args.name, tt.args.value, tt.args.project)
			if tt.wantErrorMessage != "" {
				s.Assert().Error(err)
				s.Assert().EqualError(err, tt.wantErrorMessage)
				return
			}
			s.Assert().NoError(err)

			got, err := s.client.Get(tt.args.name, tt.args.project)
			s.Assert().NoError(err)
			s.Assert().Equal(tt.want, got)
		})
	}

	s.deleteSecrets(existingSecrets, project)
}

func (s *VaultSecretStorageClientTestSuite) TestList() {
	type args struct {
		project string
	}

	project := "test-list"
	existingSecrets := map[string]string{
		"secret_1": "value_1",
		"secret_2": "value_2",
	}

	s.initializeSecrets(existingSecrets, project)

	tests := []struct {
		name             string
		args             args
		want             map[string]string
		wantErrorMessage string
	}{
		{
			name: "success: list secrets in same project as existing secrets",
			args: args{
				project: project,
			},
			want: existingSecrets,
		},
		{
			name: "success: list secrets in different project as existing secrets",
			args: args{
				project: "test-list-2",
			},
			want: make(map[string]string),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.client.List(tt.args.project)
			if tt.wantErrorMessage != "" {
				s.Assert().Error(err)
				s.Assert().EqualError(err, tt.wantErrorMessage)
				return
			}
			s.Assert().NoError(err)
			s.Assert().Equal(tt.want, got)
		})
	}

	s.deleteSecrets(existingSecrets, project)
}

func (s *VaultSecretStorageClientTestSuite) TestDelete() {
	type args struct {
		name    string
		project string
	}

	project := "test-delete"
	existingSecrets := map[string]string{
		"secret_1": "value_1",
	}

	s.initializeSecrets(existingSecrets, project)

	tests := []struct {
		name             string
		args             args
		wantErrorMessage string
	}{
		{
			name: "success: delete existing secret",
			args: args{
				name:    "secret_1",
				project: project,
			},
		},
		{
			name: "success: delete non existent secret in same project as existing secret",
			args: args{
				name:    "secret_2",
				project: project,
			},
		},
		{
			name: "success: delete non existent secret in different project as existing secret",
			args: args{
				name:    "secret_2",
				project: "test-delete-2",
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {

			err := s.client.Delete(tt.args.name, tt.args.project)
			if tt.wantErrorMessage != "" {
				s.Assert().Error(err)
				s.Assert().EqualError(err, tt.wantErrorMessage)
				return
			}
			s.Assert().NoError(err)

			_, err = s.client.Get(tt.args.name, tt.args.project)
			s.Assert().Error(err)
			s.Assert().EqualError(err, fmt.Sprintf("secret %s not found in project %s", tt.args.name, tt.args.project))
		})
	}
}

func (s *VaultSecretStorageClientTestSuite) initializeSecrets(secrets map[string]string, project string) {
	for k, v := range secrets {
		err := s.client.Set(k, v, project)
		if err != nil {
			s.FailNow("failed to create secret", err)
		}
	}
}

func (s *VaultSecretStorageClientTestSuite) deleteSecrets(secrets map[string]string, project string) {
	for k := range secrets {
		err := s.client.Delete(k, project)
		if err != nil {
			s.FailNow("failed to delete secret", err)
		}
	}
}

func TestVaultSecretStorageClient(t *testing.T) {
	suite.Run(t, new(VaultSecretStorageClientTestSuite))
}
