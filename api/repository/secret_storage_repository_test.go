//go:build integration

package repository

import (
	"fmt"
	"testing"

	"github.com/caraml-dev/mlp/api/it/database"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var gce = models.GCEGCPAuthType

type SecretStorageTestSuite struct {
	suite.Suite

	ssRepository SecretStorageRepository
	cleanupFn    func()

	testProject               *models.Project
	projectSecretStorage      *models.SecretStorage
	globalSecretStorage       *models.SecretStorage
	InternalSecretStorageType *models.SecretStorage
}

func (suite *SecretStorageTestSuite) SetupSuite() {
	db, cleanupFn, err := database.CreateTestDatabase()
	if err != nil {
		suite.FailNow(err.Error(), "Failed to connect to test database")
		return
	}

	projectRepo := NewProjectRepository(db)
	project := &models.Project{
		Name:              "test-project",
		MLFlowTrackingURL: "http://mlflow:5000",
	}
	suite.testProject, err = projectRepo.Save(project)

	suite.ssRepository = NewSecretStorageRepository(db)
	suite.cleanupFn = cleanupFn

	suite.projectSecretStorage, err = suite.ssRepository.Save(&models.SecretStorage{
		Name:      "project-secret-storage",
		Type:      models.VaultSecretStorageType,
		ProjectID: &suite.testProject.ID,
		Project:   suite.testProject,
		Scope:     models.ProjectSecretStorageScope,
		Config: models.SecretStorageConfig{
			VaultConfig: &models.VaultConfig{
				URL:         "http://vault:8200",
				KVVersion:   "2",
				AuthMethod:  models.GCPAuthMethod,
				GCPAuthType: &gce,
			},
		},
	})

	if err != nil {
		suite.Fail(err.Error(), "Failed to create project-scoped vault secret storage")
		return
	}

	suite.globalSecretStorage, err = suite.ssRepository.Save(&models.SecretStorage{
		Name:  "global-secret-storage",
		Type:  models.VaultSecretStorageType,
		Scope: models.GlobalSecretStorageScope,
		Config: models.SecretStorageConfig{
			VaultConfig: &models.VaultConfig{
				URL:         "http://vault:8200",
				KVVersion:   "2",
				AuthMethod:  models.GCPAuthMethod,
				GCPAuthType: &gce,
			},
		},
	})

	if err != nil {
		suite.Fail(err.Error(), "Failed to create global-scoped vault secret storage")
		return
	}

	suite.InternalSecretStorageType, err = suite.ssRepository.GetGlobal("internal")
	if err != nil {
		suite.Fail(err.Error(), "Failed to get internal secret storage")
		return
	}
}

func (suite *SecretStorageTestSuite) TearDownAllSuite() {
	fmt.Println("TearDownAllSuite")
	suite.cleanupFn()
}

func (suite *SecretStorageTestSuite) TestSave() {
	type args struct {
		secretStorage *models.SecretStorage
	}
	tests := []struct {
		name             string
		args             args
		want             *models.SecretStorage
		wantErrorMessage string
	}{
		{
			name: "success: create global-scoped vault secret storage",
			args: args{
				secretStorage: &models.SecretStorage{
					Name:  "global-vault-secret-storage",
					Type:  models.VaultSecretStorageType,
					Scope: models.GlobalSecretStorageScope,
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							KVVersion:   "2",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: &gce,
						},
					},
				},
			},
			want: &models.SecretStorage{
				Name:  "global-vault-secret-storage",
				Type:  models.VaultSecretStorageType,
				Scope: models.GlobalSecretStorageScope,
				Config: models.SecretStorageConfig{
					VaultConfig: &models.VaultConfig{
						URL:         "http://vault:8200",
						KVVersion:   "2",
						AuthMethod:  models.GCPAuthMethod,
						GCPAuthType: &gce,
					},
				},
			},
		},
		{
			name: "success: create project-scoped vault secret storage",
			args: args{
				secretStorage: &models.SecretStorage{
					Name:      "project-vault-secret-storage",
					Type:      models.VaultSecretStorageType,
					ProjectID: &suite.testProject.ID,
					Project:   suite.testProject,
					Scope:     models.ProjectSecretStorageScope,
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							KVVersion:   "2",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: &gce,
						},
					},
				},
			},
			want: &models.SecretStorage{
				Name:      "project-vault-secret-storage",
				Type:      models.VaultSecretStorageType,
				Scope:     models.ProjectSecretStorageScope,
				ProjectID: &suite.testProject.ID,
				Project:   suite.testProject,
				Config: models.SecretStorageConfig{
					VaultConfig: &models.VaultConfig{
						URL:         "http://vault:8200",
						KVVersion:   "2",
						AuthMethod:  models.GCPAuthMethod,
						GCPAuthType: &gce,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := suite.ssRepository.Save(tt.args.secretStorage)

			if tt.wantErrorMessage != "" {
				suite.Assert().EqualError(err, tt.wantErrorMessage)
				return
			}

			assertEqualSecretStorage(suite.T(), tt.want, got)
		})
	}
}

func (suite *SecretStorageTestSuite) TestGet() {
	type args struct {
		name    string
		project string
	}
	tests := []struct {
		name             string
		args             args
		want             *models.SecretStorage
		wantErrorMessage string
	}{
		{
			name: "success: get existing project-scoped secret storage",
			args: args{
				name:    suite.projectSecretStorage.Name,
				project: suite.testProject.Name,
			},
			want: suite.projectSecretStorage,
		},
		{
			name: "failed: not found",
			args: args{
				name:    "my-storage",
				project: suite.testProject.Name,
			},
			wantErrorMessage: "record not found",
		},
	}
	for _, tt := range tests {
		got, err := suite.ssRepository.Get(tt.args.name, tt.args.project)

		if tt.wantErrorMessage != "" {
			suite.Assert().EqualError(err, tt.wantErrorMessage)
			return
		}

		assertEqualSecretStorage(suite.T(), tt.want, got)
	}
}

func (suite *SecretStorageTestSuite) TestList() {
	type args struct {
		project string
	}
	tests := []struct {
		name             string
		args             args
		want             []*models.SecretStorage
		wantErrorMessage string
	}{
		{
			name: "success: list existing project-scoped secret storage",
			args: args{
				project: suite.testProject.Name,
			},
			want: []*models.SecretStorage{suite.projectSecretStorage},
		},
		{
			name: "no project-scoped secret storage",
			args: args{
				project: "my-project",
			},
			want: []*models.SecretStorage{},
		},
	}
	for _, tt := range tests {
		got, err := suite.ssRepository.List(tt.args.project)

		if tt.wantErrorMessage != "" {
			suite.Assert().EqualError(err, tt.wantErrorMessage)
			return
		}

		suite.Assert().Len(got, len(tt.want))
		for i, s := range got {
			assertEqualSecretStorage(suite.T(), tt.want[i], s)
		}
	}
}

func (suite *SecretStorageTestSuite) TestListGlobal() {
	tests := []struct {
		name             string
		want             []*models.SecretStorage
		wantErrorMessage string
	}{
		{
			name: "success: list existing global-scoped secret storage",
			want: []*models.SecretStorage{suite.InternalSecretStorageType, suite.globalSecretStorage},
		},
	}
	for _, tt := range tests {
		got, err := suite.ssRepository.ListGlobal()

		if tt.wantErrorMessage != "" {
			suite.Assert().EqualError(err, tt.wantErrorMessage)
			return
		}

		suite.Assert().Len(got, len(tt.want))
		for i, s := range got {
			assertEqualSecretStorage(suite.T(), tt.want[i], s)
		}
	}
}

func (suite *SecretStorageTestSuite) TestGetGlobal() {
	type args struct {
		name string
	}
	tests := []struct {
		name             string
		args             args
		want             *models.SecretStorage
		wantErrorMessage string
	}{
		{
			name: "success: get existing global-scoped secret storage",
			args: args{
				name: suite.globalSecretStorage.Name,
			},
			want: suite.globalSecretStorage,
		},
		{
			name: "failed: not found",
			args: args{
				name: "my-storage",
			},
			wantErrorMessage: "record not found",
		},
	}
	for _, tt := range tests {
		got, err := suite.ssRepository.GetGlobal(tt.args.name)

		if tt.wantErrorMessage != "" {
			suite.Assert().EqualError(err, tt.wantErrorMessage)
			return
		}

		assertEqualSecretStorage(suite.T(), tt.want, got)
	}
}

func (suite *SecretStorageTestSuite) TestDelete() {
	type args struct {
		secretStorage *models.SecretStorage
	}
	tests := []struct {
		name             string
		args             args
		wantErrorMessage string
	}{
		{
			name: "success: delete project-scoped vault secret storage",
			args: args{
				secretStorage: &models.SecretStorage{
					Name:      "project-vault-secret-storage",
					Type:      models.VaultSecretStorageType,
					ProjectID: &suite.testProject.ID,
					Project:   suite.testProject,
					Scope:     models.ProjectSecretStorageScope,
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							KVVersion:   "2",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: &gce,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := suite.ssRepository.Save(tt.args.secretStorage)

			if tt.wantErrorMessage != "" {
				suite.Assert().EqualError(err, tt.wantErrorMessage)
				return
			}

			assertEqualSecretStorage(suite.T(), tt.args.secretStorage, got)

			err = suite.ssRepository.Delete(got)
			suite.Assert().NoError(err)

			_, err = suite.ssRepository.Get(got.Name, got.Project.Name)
			suite.Assert().EqualError(err, "record not found")
		})
	}
}

func TestSecretStorageRepository(t *testing.T) {
	suite.Run(t, new(SecretStorageTestSuite))
}

func assertEqualSecretStorage(t *testing.T, exp *models.SecretStorage, got *models.SecretStorage) {
	assert.Equal(t, exp.Name, got.Name)
	assert.Equal(t, exp.Type, got.Type)
	assert.Equal(t, exp.Scope, got.Scope)
	assert.Equal(t, exp.Config, got.Config)
	if exp.Scope == models.ProjectSecretStorageScope {
		assert.Equal(t, exp.ProjectID, got.ProjectID)
		assert.Equal(t, exp.Project.Name, got.Project.Name)
	}
}
