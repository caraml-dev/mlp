package service

import (
	"fmt"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/secretstorage"
	ssmocks "github.com/caraml-dev/mlp/api/pkg/secretstorage/mocks"
	"github.com/caraml-dev/mlp/api/repository/mocks"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSecretService_FindByID(t *testing.T) {
	internalSecretStorage := &models.SecretStorage{
		ID:   1,
		Name: "internal-secret-storage",
		Type: models.InternalSecretStorageType,
	}

	vaultSecretStorage := &models.SecretStorage{
		ID:   2,
		Name: "vault-secret-storage",
		Type: models.VaultSecretStorageType,
	}

	project := &models.Project{
		ID:   models.ID(1),
		Name: "project",
	}

	type args struct {
		secretID models.ID
	}

	tests := []struct {
		name                             string
		args                             args
		existingSecret                   *models.Secret
		errorFromSecretRepository        error
		errorFromSecretStorageRepository error
		errorFromSecretStorageClient     error
		expectedError                    string
	}{
		{
			name: "success: get existing internal secret",
			args: args{
				secretID: models.ID(1),
			},
			existingSecret: &models.Secret{
				ID:              models.ID(1),
				ProjectID:       project.ID,
				Project:         project,
				Name:            "name",
				Data:            "plainData",
				SecretStorageID: &internalSecretStorage.ID,
				SecretStorage:   internalSecretStorage,
			},
		},
		{
			name: "success: get existing external secret",
			args: args{
				secretID: models.ID(1),
			},
			existingSecret: &models.Secret{
				ID:              models.ID(1),
				ProjectID:       project.ID,
				Project:         project,
				Name:            "name",
				Data:            "plainData",
				SecretStorageID: &vaultSecretStorage.ID,
				SecretStorage:   vaultSecretStorage,
			},
		},
		{
			name: "error: secret not found",
			args: args{
				secretID: models.ID(1),
			},
			existingSecret:            nil,
			errorFromSecretRepository: gorm.ErrRecordNotFound,
			expectedError:             "error when fetching secret with id: 1, error: record not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRepository := &mocks.ProjectRepository{}
			if tt.existingSecret != nil {
				projectRepository.On("Get", tt.existingSecret.ProjectID).Return(project, nil)
			}

			ssClientRegistry, err := secretstorage.NewRegistry([]*models.SecretStorage{})
			require.NoError(t, err)

			ssClient := &ssmocks.Client{}
			if tt.existingSecret != nil {
				ssClient.On("Get", tt.existingSecret.Name, project.Name).Return(tt.existingSecret.Data, tt.errorFromSecretStorageClient)
			}
			ssClientRegistry.Set(internalSecretStorage.ID, ssClient)
			ssClientRegistry.Set(vaultSecretStorage.ID, ssClient)

			storageRepository := &mocks.SecretStorageRepository{}
			storageRepository.On("Get", internalSecretStorage.ID).Return(internalSecretStorage, tt.errorFromSecretStorageRepository)
			storageRepository.On("Get", vaultSecretStorage.ID).Return(vaultSecretStorage, tt.errorFromSecretStorageRepository)

			secretRepository := &mocks.SecretRepository{}
			secretRepository.On("Get", tt.args.secretID).Return(tt.existingSecret, tt.errorFromSecretRepository)

			secretService := NewSecretService(secretRepository, storageRepository, projectRepository, ssClientRegistry, vaultSecretStorage)
			result, err := secretService.FindByID(tt.args.secretID)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.existingSecret.ID, result.ID)
			assert.Equal(t, tt.existingSecret.Name, result.Name)
			assert.Equal(t, tt.existingSecret.ProjectID, result.ProjectID)
			assert.Equal(t, tt.existingSecret.Data, result.Data)
			if tt.existingSecret.SecretStorageID != nil {
				assert.Equal(t, *tt.existingSecret.SecretStorageID, *result.SecretStorageID)
			} else {
				assert.Equal(t, vaultSecretStorage.ID, *result.SecretStorageID)
			}
			require.NoError(t, err)
		})
	}
}

func TestSecretService_Create(t *testing.T) {
	internalSecretStorage := &models.SecretStorage{
		ID:   1,
		Name: "internal-secret-storage",
		Type: models.InternalSecretStorageType,
	}

	vaultSecretStorage := &models.SecretStorage{
		ID:   2,
		Name: "vault-secret-storage",
		Type: models.VaultSecretStorageType,
	}

	project := &models.Project{
		ID:   models.ID(1),
		Name: "project",
	}

	tests := []struct {
		name                             string
		secret                           *models.Secret
		errorFromSecretRepository        error
		errorFromSecretStorageRepository error
		errorFromSecretStorageClient     error
		expectedError                    string
	}{
		{
			name: "success: using default storage",
			secret: &models.Secret{
				ID:        models.ID(1),
				ProjectID: project.ID,
				Name:      "name",
				Data:      "plainData",
			},
		},
		{
			name: "success: using internal storage",
			secret: &models.Secret{
				ID:              models.ID(1),
				ProjectID:       project.ID,
				SecretStorageID: &internalSecretStorage.ID,
				Name:            "name",
				Data:            "plainData",
			},
		},
		{
			name: "success: using vault storage",
			secret: &models.Secret{
				ID:              models.ID(1),
				ProjectID:       project.ID,
				SecretStorageID: &vaultSecretStorage.ID,
				Name:            "name",
				Data:            "plainData",
			},
		},
		{
			name: "error: should raise error when failed save to db",
			secret: &models.Secret{
				ID:        models.ID(1),
				ProjectID: models.ID(1),
				Name:      "name",
				Data:      "plainData",
			},
			errorFromSecretRepository: fmt.Errorf("db is down"),
			expectedError:             "error when saving secret in database, error: db is down",
		},
		{
			name: "error: should raise error when failed to store secret",
			secret: &models.Secret{
				ID:              models.ID(1),
				ProjectID:       models.ID(1),
				SecretStorageID: &vaultSecretStorage.ID,
				Name:            "name",
				Data:            "plainData",
			},
			errorFromSecretStorageClient: fmt.Errorf("vault is down"),
			expectedError:                "error when creating secret in secret storage with id: 2, error: vault is down",
		},
		{
			name: "error: should raise error when secret storage is not found",
			secret: &models.Secret{
				ID:              models.ID(1),
				ProjectID:       models.ID(1),
				SecretStorageID: &vaultSecretStorage.ID,
				Name:            "name",
				Data:            "plainData",
			},
			errorFromSecretStorageRepository: fmt.Errorf("secret storage not found"),
			expectedError:                    "error when fetching secret storage with id: 2, error: secret storage not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRepository := &mocks.ProjectRepository{}
			projectRepository.On("Get", tt.secret.ProjectID).Return(project, nil)

			ssClientRegistry, err := secretstorage.NewRegistry([]*models.SecretStorage{})
			require.NoError(t, err)

			ssClient := &ssmocks.Client{}
			ssClient.On("Set", tt.secret.Name, tt.secret.Data, project.Name).Return(tt.errorFromSecretStorageClient)
			ssClientRegistry.Set(internalSecretStorage.ID, ssClient)
			ssClientRegistry.Set(vaultSecretStorage.ID, ssClient)

			storageRepository := &mocks.SecretStorageRepository{}
			storageRepository.On("Get", internalSecretStorage.ID).Return(internalSecretStorage, tt.errorFromSecretStorageRepository)
			storageRepository.On("Get", vaultSecretStorage.ID).Return(vaultSecretStorage, tt.errorFromSecretStorageRepository)

			secretRepository := &mocks.SecretRepository{}
			secretRepository.On("Save", tt.secret).Return(tt.secret, tt.errorFromSecretRepository)

			secretService := NewSecretService(secretRepository, storageRepository, projectRepository, ssClientRegistry, vaultSecretStorage)
			result, err := secretService.Create(tt.secret)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.secret.ID, result.ID)
			assert.Equal(t, tt.secret.Name, result.Name)
			assert.Equal(t, tt.secret.ProjectID, result.ProjectID)
			assert.Equal(t, tt.secret.Data, result.Data)
			if tt.secret.SecretStorageID != nil {
				assert.Equal(t, *tt.secret.SecretStorageID, *result.SecretStorageID)
			} else {
				assert.Equal(t, vaultSecretStorage.ID, *result.SecretStorageID)
			}
			require.NoError(t, err)
		})
	}
}

func TestSecretService_Delete(t *testing.T) {
	internalSecretStorage := &models.SecretStorage{
		ID:   1,
		Name: "internal-secret-storage",
		Type: models.InternalSecretStorageType,
	}

	vaultSecretStorage := &models.SecretStorage{
		ID:   2,
		Name: "vault-secret-storage",
		Type: models.VaultSecretStorageType,
	}

	project := &models.Project{
		ID:   models.ID(1),
		Name: "project",
	}

	tests := []struct {
		name                         string
		secretID                     models.ID
		existingSecret               *models.Secret
		errorFromSecretRepository    error
		errorFromSecretStorageClient error
		expectedError                string
	}{
		{
			name:     "success: delete internal secret",
			secretID: models.ID(1),
			existingSecret: &models.Secret{
				ID:              models.ID(1),
				Name:            "my-secret",
				ProjectID:       project.ID,
				Project:         project,
				SecretStorageID: &internalSecretStorage.ID,
				SecretStorage:   internalSecretStorage,
			},
		},
		{
			name:     "success: delete external secret",
			secretID: models.ID(1),
			existingSecret: &models.Secret{
				ID:              models.ID(1),
				Name:            "my-secret",
				ProjectID:       project.ID,
				Project:         project,
				SecretStorageID: &vaultSecretStorage.ID,
				SecretStorage:   vaultSecretStorage,
			},
		},
		{
			name:                      "error: should return error when failed to delete secret",
			secretID:                  models.ID(1),
			errorFromSecretRepository: fmt.Errorf("db is down"),
			expectedError:             "error when fetching secret with id: 1, error: db is down",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretRepository := &mocks.SecretRepository{}
			secretRepository.On("Get", tt.secretID).Return(tt.existingSecret, tt.errorFromSecretRepository)
			secretRepository.On("Delete", tt.secretID).Return(tt.errorFromSecretRepository)

			ssClientRegistry, err := secretstorage.NewRegistry([]*models.SecretStorage{})
			require.NoError(t, err)
			ssClient := &ssmocks.Client{}

			if tt.existingSecret != nil {
				ssClient.On("Delete", tt.existingSecret.Name, project.Name).Return(tt.errorFromSecretStorageClient)
			}
			ssClientRegistry.Set(internalSecretStorage.ID, ssClient)
			ssClientRegistry.Set(vaultSecretStorage.ID, ssClient)

			secretService := NewSecretService(secretRepository, nil, nil, ssClientRegistry, vaultSecretStorage)

			err = secretService.Delete(tt.secretID)
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestSecretService_List(t *testing.T) {
	internalSecretStorage := &models.SecretStorage{
		ID:   1,
		Name: "internal-secret-storage",
		Type: models.InternalSecretStorageType,
	}

	vaultSecretStorage := &models.SecretStorage{
		ID:   2,
		Name: "vault-secret-storage",
		Type: models.VaultSecretStorageType,
	}

	project := &models.Project{
		ID:   models.ID(1),
		Name: "project",
	}

	secrets := []*models.Secret{
		{
			ID:              models.ID(1),
			ProjectID:       project.ID,
			Project:         project,
			SecretStorageID: &internalSecretStorage.ID,
			SecretStorage:   internalSecretStorage,
			Name:            "name1",
			Data:            "plainData",
		},
		{
			ID:              models.ID(2),
			ProjectID:       project.ID,
			Project:         project,
			SecretStorageID: &vaultSecretStorage.ID,
			SecretStorage:   vaultSecretStorage,
			Name:            "name2",
			Data:            "plainData",
		},
		{
			ID:              models.ID(3),
			ProjectID:       project.ID,
			Project:         project,
			SecretStorageID: &vaultSecretStorage.ID,
			SecretStorage:   vaultSecretStorage,
			Name:            "name3",
			Data:            "plainData3",
		},
	}

	projectRepository := &mocks.ProjectRepository{}
	projectRepository.On("Get", project.ID).Return(project, nil)

	secretRepository := &mocks.SecretRepository{}
	secretRepository.On("List", project.ID).Return(secrets, nil)

	ssClientRegistry, err := secretstorage.NewRegistry([]*models.SecretStorage{})
	require.NoError(t, err)

	ssClient := &ssmocks.Client{}
	ssClient.On("List", project.Name).Return(map[string]string{
		secrets[0].Name: secrets[0].Data,
	}, nil)
	ssClient.On("List", project.Name).Return(map[string]string{
		secrets[1].Name: secrets[1].Data,
		secrets[2].Name: secrets[2].Data,
	}, nil)
	ssClientRegistry.Set(internalSecretStorage.ID, ssClient)
	ssClientRegistry.Set(vaultSecretStorage.ID, ssClient)

	storageRepository := &mocks.SecretStorageRepository{}

	secretService := NewSecretService(secretRepository, storageRepository, projectRepository, ssClientRegistry, vaultSecretStorage)
	actual, err := secretService.List(project.ID)
	assert.NoError(t, err)
	assert.Equal(t, secrets, actual)
}

func TestSecretService_Update(t *testing.T) {
	internalSecretStorage := &models.SecretStorage{
		ID:   1,
		Name: "internal-secret-storage",
		Type: models.InternalSecretStorageType,
	}

	vaultSecretStorage := &models.SecretStorage{
		ID:   2,
		Name: "vault-secret-storage",
		Type: models.VaultSecretStorageType,
	}

	project := &models.Project{
		ID:   models.ID(1),
		Name: "project",
	}

	existingSecret := &models.Secret{
		ID:              models.ID(1),
		ProjectID:       project.ID,
		Project:         project,
		SecretStorageID: &internalSecretStorage.ID,
		SecretStorage:   internalSecretStorage,
		Name:            "name1",
		Data:            "plainData",
	}

	type args struct {
		secret *models.Secret
	}

	tests := []struct {
		name          string
		args          args
		want          *models.Secret
		expectedError string
	}{
		{
			name: "success: update value",
			args: args{
				secret: &models.Secret{
					ID:              existingSecret.ID,
					ProjectID:       project.ID,
					Project:         project,
					SecretStorageID: &internalSecretStorage.ID,
					SecretStorage:   internalSecretStorage,
					Name:            "name1",
					Data:            "plainData2",
				},
			},
			want: &models.Secret{
				ID:              existingSecret.ID,
				ProjectID:       project.ID,
				Project:         project,
				SecretStorageID: &internalSecretStorage.ID,
				SecretStorage:   internalSecretStorage,
				Name:            "name1",
				Data:            "plainData2",
			},
		},
		{
			name: "success: migrate storage",
			args: args{
				secret: &models.Secret{
					ID:              existingSecret.ID,
					ProjectID:       project.ID,
					Project:         project,
					SecretStorageID: &vaultSecretStorage.ID,
					SecretStorage:   vaultSecretStorage,
					Name:            "name1",
				},
			},
			want: &models.Secret{
				ID:              existingSecret.ID,
				ProjectID:       project.ID,
				Project:         project,
				SecretStorageID: &vaultSecretStorage.ID,
				SecretStorage:   vaultSecretStorage,
				Name:            "name1",
				Data:            "plainData",
			},
		},
		{
			name: "success: migrate storage and update secret",
			args: args{
				secret: &models.Secret{
					ID:              existingSecret.ID,
					ProjectID:       project.ID,
					Project:         project,
					SecretStorageID: &vaultSecretStorage.ID,
					SecretStorage:   vaultSecretStorage,
					Name:            "name1",
					Data:            "plainData2",
				},
			},
			want: &models.Secret{
				ID:              existingSecret.ID,
				ProjectID:       project.ID,
				Project:         project,
				SecretStorageID: &vaultSecretStorage.ID,
				SecretStorage:   vaultSecretStorage,
				Name:            "name1",
				Data:            "plainData2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			projectRepository := &mocks.ProjectRepository{}
			projectRepository.On("Get", project.ID).Return(project, nil)

			secretRepository := &mocks.SecretRepository{}
			secretRepository.On("Get", existingSecret.ID).Return(existingSecret, nil)
			secretRepository.On("Save", tt.args.secret).Return(tt.args.secret, nil)

			ssClientRegistry, err := secretstorage.NewRegistry([]*models.SecretStorage{})
			require.NoError(t, err)

			ssClient := &ssmocks.Client{}
			ssClient.On("Get", existingSecret.Name, project.Name).Return(existingSecret.Data, nil)
			ssClient.On("Set", tt.args.secret.Name, tt.args.secret.Data, project.Name).Return(nil)
			ssClient.On("Set", tt.args.secret.Name, existingSecret.Data, project.Name).Return(nil)
			ssClient.On("Delete", tt.args.secret.Name, project.Name).Return(nil)

			ssClientRegistry.Set(internalSecretStorage.ID, ssClient)
			ssClientRegistry.Set(vaultSecretStorage.ID, ssClient)

			storageRepository := &mocks.SecretStorageRepository{}
			storageRepository.On("Get", internalSecretStorage.ID).Return(internalSecretStorage, nil)
			storageRepository.On("Get", vaultSecretStorage.ID).Return(vaultSecretStorage, nil)

			secretService := NewSecretService(secretRepository, storageRepository, projectRepository, ssClientRegistry, vaultSecretStorage)
			got, err := secretService.Update(tt.args.secret)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				return
			}
			assert.Equalf(t, tt.want, got, "Update(%v)", tt.args.secret)
		})
	}
}
