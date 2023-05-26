//go:build integration

package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/caraml-dev/mlp/api/it/database"
	"github.com/caraml-dev/mlp/api/models"
)

type SecretRepositoryTestSuite struct {
	suite.Suite
	secretRepository SecretRepository
	cleanupFn        func()

	project               *models.Project
	internalSecretStorage *models.SecretStorage
	existingSecrets       []*models.Secret
}

func (suite *SecretRepositoryTestSuite) SetupSuite() {
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

	ssRepo := NewSecretStorageRepository(db)
	suite.internalSecretStorage, err = ssRepo.GetGlobal("internal")
	suite.Require().NoError(err, "Failed to get internal secret storage")

	suite.project, err = projectRepo.Save(project)
	suite.Require().NoError(err, "Failed to create project")

	suite.secretRepository = NewSecretRepository(db)
	suite.cleanupFn = cleanupFn

	suite.existingSecrets = make([]*models.Secret, 0)
	for i := 0; i < 3; i++ {
		s, err := suite.secretRepository.Save(&models.Secret{
			ProjectID:       suite.project.ID,
			SecretStorageID: &suite.internalSecretStorage.ID,
			Name:            fmt.Sprintf("secret-%d", i),
			Data:            fmt.Sprintf("data-%d", i),
		})

		suite.Require().NoError(err, "Failed to create secret")
		suite.existingSecrets = append(suite.existingSecrets, s)
	}
}

func (suite *SecretRepositoryTestSuite) TearDownSuite() {
	suite.cleanupFn()
}

func (suite *SecretRepositoryTestSuite) TestSave() {
	tests := []struct {
		name       string
		secret     *models.Secret
		want       *models.Secret
		wantErrMsg string
	}{
		{
			name: "Should success if all validation is met",
			secret: &models.Secret{
				ProjectID:       models.ID(1),
				SecretStorageID: &suite.internalSecretStorage.ID,
				Name:            "secret_name",
				Data:            "data",
			},
			want: &models.Secret{
				ProjectID:       models.ID(1),
				SecretStorageID: &suite.internalSecretStorage.ID,
				Name:            "secret_name",
				Data:            "data",
			},
		},
		{
			name: "Should failed if project_id is not exist in db",
			secret: &models.Secret{
				ProjectID: models.ID(2),
				Name:      "name",
				Data:      "data",
			},
			wantErrMsg: `pq: insert or update on table "secrets" violates foreign key constraint "secrets_project_id_fkey"`,
		},
		{
			name: "Should failed if existing secret name used in the same project_id",
			secret: &models.Secret{
				ProjectID: models.ID(1),
				Name:      "secret_name",
				Data:      "data",
			},
			wantErrMsg: `pq: duplicate key value violates unique constraint "secrets_project_id_name_key"`,
		},
		{
			name: "Should success edit secret data",
			secret: &models.Secret{
				ID:              suite.existingSecrets[0].ID,
				ProjectID:       suite.existingSecrets[0].ProjectID,
				SecretStorageID: suite.existingSecrets[0].SecretStorageID,
				Name:            suite.existingSecrets[0].Name,
				Data:            "new_data",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: time.Now(),
				},
			},
			want: &models.Secret{
				ID:              suite.existingSecrets[0].ID,
				ProjectID:       suite.existingSecrets[0].ProjectID,
				SecretStorageID: suite.existingSecrets[0].SecretStorageID,
				Name:            suite.existingSecrets[0].Name,
				Data:            "new_data",
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := suite.secretRepository.Save(tt.secret)
			if tt.wantErrMsg != "" {
				suite.Assert().EqualError(err, tt.wantErrMsg)
				return
			}

			suite.Require().NoError(err)
			assertSecretEqual(suite.T(), tt.want, got)
		})
	}
}

func (suite *SecretRepositoryTestSuite) TestDelete() {
	secret, err := suite.secretRepository.Save(&models.Secret{
		Name:            "will-be-deleted",
		ProjectID:       suite.project.ID,
		SecretStorageID: &suite.internalSecretStorage.ID,
		Data:            "data",
	})

	suite.Require().NoError(err, "Failed to create secret")

	type args struct {
		id models.ID
	}

	tests := []struct {
		name             string
		args             args
		wantErrorMessage string
	}{
		{
			name: "Should success deleted secret",
			args: args{
				id: secret.ID,
			},
		},
		{
			name: "Should success even when secret not exist",
			args: args{
				id: 123,
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.secretRepository.Delete(tt.args.id)
			if tt.wantErrorMessage != "" {
				suite.Assert().EqualError(err, tt.wantErrorMessage)
				return
			}
			suite.Require().NoError(err)
		})
	}
}

func (suite *SecretRepositoryTestSuite) TestGet() {
	type args struct {
		id models.ID
	}

	tests := []struct {
		name             string
		args             args
		want             *models.Secret
		wantErrorMessage string
	}{
		{
			name: "Should success get secret",
			args: args{
				id: suite.existingSecrets[0].ID,
			},
			want: &models.Secret{
				ID:              suite.existingSecrets[0].ID,
				ProjectID:       suite.existingSecrets[0].ProjectID,
				SecretStorageID: suite.existingSecrets[0].SecretStorageID,
				Name:            suite.existingSecrets[0].Name,
				Data:            suite.existingSecrets[0].Data,
			},
		},
		{
			name: "Should failed when secret not found",
			args: args{
				id: 123,
			},
			wantErrorMessage: "secret with id 123 not found",
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			secret, err := suite.secretRepository.Get(tt.args.id)
			if tt.wantErrorMessage != "" {
				suite.Assert().EqualError(err, tt.wantErrorMessage)
				return
			}

			suite.Require().NoError(err)
			assertSecretEqual(suite.T(), tt.want, secret)
		})
	}
}

func (suite *SecretRepositoryTestSuite) TestList() {
	secrets, err := suite.secretRepository.List(suite.project.ID)
	suite.Require().NoError(err)

	for i, got := range secrets {
		assertSecretEqual(suite.T(), suite.existingSecrets[i], got)
	}
}

func assertSecretEqual(t *testing.T, want *models.Secret, got *models.Secret) {
	assert.NotZero(t, got.ID)
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Data, got.Data)
	assert.Equal(t, want.ProjectID, got.ProjectID)
	assert.Equal(t, *want.SecretStorageID, *got.SecretStorageID)

	assert.NotEmpty(t, got.CreatedAt)
	assert.NotEmpty(t, got.UpdatedAt)
}

func TestSecretRepository(t *testing.T) {
	suite.Run(t, new(SecretRepositoryTestSuite))
}
