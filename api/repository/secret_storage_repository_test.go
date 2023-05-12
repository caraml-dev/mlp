//go:build integration

package repository

import (
	"testing"

	"github.com/caraml-dev/mlp/api/it/database"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SecretStorageTestSuite struct {
	suite.Suite

	ssRepository SecretStorageRepository
	cleanupFn    func()

	project               *models.Project
	projectSecretStorage  *models.SecretStorage
	globalSecretStorage   *models.SecretStorage
	internalSecretStorage *models.SecretStorage
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
	suite.project, err = projectRepo.Save(project)
	if err != nil {
		suite.Fail(err.Error(), "Failed to create project")
		return
	}

	suite.ssRepository = NewSecretStorageRepository(db)
	suite.cleanupFn = cleanupFn

	suite.projectSecretStorage, err = suite.ssRepository.Save(&models.SecretStorage{
		Name:      "project-secret-storage",
		Type:      models.VaultSecretStorageType,
		ProjectID: &suite.project.ID,
		Project:   suite.project,
		Scope:     models.ProjectSecretStorageScope,
		Config: models.SecretStorageConfig{
			VaultConfig: &models.VaultConfig{
				URL:         "http://vault:8200",
				AuthMethod:  models.GCPAuthMethod,
				GCPAuthType: models.GCEGCPAuthType,
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
				AuthMethod:  models.GCPAuthMethod,
				GCPAuthType: models.GCEGCPAuthType,
			},
		},
	})

	if err != nil {
		suite.Fail(err.Error(), "Failed to create global-scoped vault secret storage")
		return
	}

	suite.internalSecretStorage, err = suite.ssRepository.GetGlobal("internal")
	if err != nil {
		suite.Fail(err.Error(), "Failed to get internal secret storage")
		return
	}
}

func (suite *SecretStorageTestSuite) TearDownAllSuite() {
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
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
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
						AuthMethod:  models.GCPAuthMethod,
						GCPAuthType: models.GCEGCPAuthType,
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
					ProjectID: &suite.project.ID,
					Project:   suite.project,
					Scope:     models.ProjectSecretStorageScope,
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
						},
					},
				},
			},
			want: &models.SecretStorage{
				Name:      "project-vault-secret-storage",
				Type:      models.VaultSecretStorageType,
				Scope:     models.ProjectSecretStorageScope,
				ProjectID: &suite.project.ID,
				Project:   suite.project,
				Config: models.SecretStorageConfig{
					VaultConfig: &models.VaultConfig{
						URL:         "http://vault:8200",
						AuthMethod:  models.GCPAuthMethod,
						GCPAuthType: models.GCEGCPAuthType,
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
		ID models.ID
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
				ID: suite.projectSecretStorage.ID,
			},
			want: suite.projectSecretStorage,
		},
		{
			name: "failed: not found",
			args: args{
				ID: 100,
			},
			wantErrorMessage: "record not found",
		},
	}
	for _, tt := range tests {
		got, err := suite.ssRepository.Get(tt.args.ID)

		if tt.wantErrorMessage != "" {
			suite.Assert().EqualError(err, tt.wantErrorMessage)
			return
		}

		assertEqualSecretStorage(suite.T(), tt.want, got)
	}
}

func (suite *SecretStorageTestSuite) TestList() {
	type args struct {
		projectID models.ID
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
				projectID: suite.project.ID,
			},
			want: []*models.SecretStorage{suite.projectSecretStorage},
		},
		{
			name: "no project-scoped secret storage",
			args: args{
				projectID: 10,
			},
			want: []*models.SecretStorage{},
		},
	}
	for _, tt := range tests {
		got, err := suite.ssRepository.List(tt.args.projectID)

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
	got, err := suite.ssRepository.ListGlobal()
	suite.Assert().NoError(err)
	exps := []*models.SecretStorage{
		suite.internalSecretStorage, suite.globalSecretStorage,
	}

	suite.Assert().Len(got, len(exps))
	for i, exp := range exps {
		assertEqualSecretStorage(suite.T(), exp, got[i])
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
					ProjectID: &suite.project.ID,
					Project:   suite.project,
					Scope:     models.ProjectSecretStorageScope,
					Config: models.SecretStorageConfig{
						VaultConfig: &models.VaultConfig{
							URL:         "http://vault:8200",
							AuthMethod:  models.GCPAuthMethod,
							GCPAuthType: models.GCEGCPAuthType,
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

			err = suite.ssRepository.Delete(got.ID)
			suite.Assert().NoError(err)

			_, err = suite.ssRepository.Get(got.ID)
			suite.Assert().EqualError(err, "record not found")
		})
	}
}

func (suite *SecretStorageTestSuite) TestListAll() {
	got, err := suite.ssRepository.ListAll()
	suite.Assert().NoError(err)
	exps := []*models.SecretStorage{
		suite.internalSecretStorage, suite.projectSecretStorage, suite.globalSecretStorage,
	}

	suite.Assert().Len(got, len(exps))
	for i, exp := range exps {
		assertEqualSecretStorage(suite.T(), exp, got[i])
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
