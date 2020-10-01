// +build integration integration_local

package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gojek/mlp/api/it/database"
	"github.com/gojek/mlp/api/models"
	"github.com/gojek/mlp/api/util"
)

const (
	passphrase = "password"
)

func TestSave(t *testing.T) {
	testCases := []struct {
		desc           string
		secret         *models.Secret
		existingSecret *models.Secret
		expectedSecret *models.Secret
		expectedError  string
	}{
		{
			desc: "Should success if all validation is met",
			secret: &models.Secret{
				ProjectId: models.Id(1),
				Name:      "secret_name",
				Data:      "data",
			},
			expectedSecret: &models.Secret{
				ProjectId: models.Id(1),
				Name:      "secret_name",
				Data:      "data",
			},
		},
		{
			desc: "Should failed if project_id is not exist in db",
			secret: &models.Secret{
				ProjectId: models.Id(2),
				Name:      "name",
				Data:      "data",
			},
			expectedError: `pq: insert or update on table "secrets" violates foreign key constraint "secrets_project_id_fkey"`,
		},
		{
			desc: "Should failed if existing secret name used in the same project_id",
			secret: &models.Secret{
				ProjectId: models.Id(1),
				Name:      "secret_name",
				Data:      "data",
			},
			existingSecret: &models.Secret{
				ProjectId: models.Id(1),
				Name:      "secret_name",
				Data:      "old_data",
			},
			expectedError: `pq: duplicate key value violates unique constraint "secrets_project_id_name_key"`,
		},
		{
			desc: "Should success edit secret data",
			secret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "secret_name",
				Data:      "data",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: time.Now(),
				},
			},
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "secret_name",
				Data:      "old_data",
			},
			expectedSecret: &models.Secret{
				ProjectId: models.Id(1),
				Name:      "secret_name",
				Data:      "data",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				projectStorage := NewProjectStorage(db)
				projectStorage.Save(&models.Project{
					Id:   models.Id(1),
					Name: "project_name",
				})
				secretStorage := NewSecretStorage(db, passphrase)
				if tC.existingSecret != nil {
					_, err := secretStorage.Save(tC.existingSecret)
					require.NoError(t, err)
				}
				result, err := secretStorage.Save(tC.secret)
				if tC.expectedError != "" {
					assert.EqualError(t, err, tC.expectedError)
				} else {
					fmt.Printf("result %+v", *result)
					require.NoError(t, err)
					assert.NotZero(t, result.Id)
					assert.NotEmpty(t, result.CreatedAt)
					assert.NotEmpty(t, result.UpdatedAt)
					assert.Equal(t, tC.expectedSecret.ProjectId, result.ProjectId)
					assert.Equal(t, tC.expectedSecret.Name, result.Name)

					plain, err := result.DecryptData(util.CreateHash(passphrase))
					assert.NoError(t, err)
					assert.Equal(t, tC.expectedSecret.Data, plain.Data)
				}
			})
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		desc           string
		existingSecret *models.Secret
		secretToDelete *models.Secret
		expectedError  string
	}{
		{
			desc: "Should success deleted secret",
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				Name:      "name",
				ProjectId: models.Id(1),
				Data:      "data",
			},
			secretToDelete: &models.Secret{
				Id:        models.Id(1),
				Name:      "name",
				ProjectId: models.Id(1),
				Data:      "data",
			},
		},
		{
			desc: "Should success even when secret not exist",
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				Name:      "name",
				ProjectId: models.Id(1),
				Data:      "data",
			},
			secretToDelete: &models.Secret{
				Id:        models.Id(2),
				Name:      "name_2",
				ProjectId: models.Id(1),
				Data:      "data",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				projectStorage := NewProjectStorage(db)
				projectStorage.Save(&models.Project{
					Id:   models.Id(1),
					Name: "project_name",
				})
				secretStorage := NewSecretStorage(db, "password")
				if tC.existingSecret != nil {
					_, err := secretStorage.Save(tC.existingSecret)
					require.NoError(t, err)
				}
				err := secretStorage.Delete(tC.secretToDelete.Id, tC.secretToDelete.ProjectId)
				if tC.expectedError != "" {
					assert.EqualError(t, err, tC.expectedError)
				} else {
					require.NoError(t, err)
				}
			})
		})
	}
}

func TestGetAsPlainText(t *testing.T) {
	testCases := []struct {
		desc           string
		existingSecret *models.Secret
		secretId       models.Id
		projectId      models.Id
		expectedSecret *models.Secret
		expectedError  string
	}{
		{
			desc: "Should success get secret",
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "data",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			secretId:  models.Id(1),
			projectId: models.Id(1),
			expectedSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "data",
			},
		},
		{
			desc: "Should failed when secret not found",
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "data",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			secretId:      models.Id(1),
			projectId:     models.Id(2),
			expectedError: "record not found",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				projectStorage := NewProjectStorage(db)
				projectStorage.Save(&models.Project{
					Id:   models.Id(1),
					Name: "project_name",
				})
				secretStorage := NewSecretStorage(db, "password")
				if tC.existingSecret != nil {
					_, err := secretStorage.Save(tC.existingSecret)
					require.NoError(t, err)
				}
				secret, err := secretStorage.GetAsPlainText(tC.secretId, tC.projectId)
				if tC.expectedError != "" {
					assert.EqualError(t, err, tC.expectedError)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tC.expectedSecret.Id, secret.Id)
					assert.Equal(t, tC.expectedSecret.Name, secret.Name)
					assert.Equal(t, tC.expectedSecret.Data, secret.Data)
					assert.Equal(t, tC.expectedSecret.ProjectId, secret.ProjectId)
				}
			})
		})
	}
}

func TestGetByNameAsPlainText(t *testing.T) {
	testCases := []struct {
		desc           string
		existingSecret *models.Secret
		name           string
		projectId      models.Id
		expectedSecret *models.Secret
		expectedError  string
	}{
		{
			desc: "Should success get secret",
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "data",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			name:      "name",
			projectId: models.Id(1),
			expectedSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "data",
			},
		},
		{
			desc: "Should failed when secret not found",
			existingSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "data",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			name:          "other-name",
			projectId:     models.Id(1),
			expectedError: "record not found",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				projectStorage := NewProjectStorage(db)
				projectStorage.Save(&models.Project{
					Id:   models.Id(1),
					Name: "project_name",
				})
				secretStorage := NewSecretStorage(db, "password")
				if tC.existingSecret != nil {
					_, err := secretStorage.Save(tC.existingSecret)
					require.NoError(t, err)
				}
				secret, err := secretStorage.GetByNameAsPlainText(tC.name, tC.projectId)
				if tC.expectedError != "" {
					assert.EqualError(t, err, tC.expectedError)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tC.expectedSecret.Id, secret.Id)
					assert.Equal(t, tC.expectedSecret.Name, secret.Name)
					assert.Equal(t, tC.expectedSecret.Data, secret.Data)
					assert.Equal(t, tC.expectedSecret.ProjectId, secret.ProjectId)
				}
			})
		})
	}
}

func TestList(t *testing.T) {
	database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
		projectStorage := NewProjectStorage(db)
		projectStorage.Save(&models.Project{
			Id:   models.Id(1),
			Name: "project_name",
		})
		secretStorage := NewSecretStorage(db, passphrase)
		secret1 := &models.Secret{
			ProjectId: 1,
			Name:      "secret-1",
			Data:      "data-1",
		}
		_, err := secretStorage.Save(secret1)
		assert.NoError(t, err)

		secret2 := &models.Secret{
			ProjectId: 1,
			Name:      "secret-2",
			Data:      "data-2",
		}
		_, err = secretStorage.Save(secret2)
		assert.NoError(t, err)

		secrets, err := secretStorage.List(1)
		assert.NoError(t, err)
		assert.Len(t, secrets, 2)

		// assert both secret is encrypted
		plainSecret1, err := secrets[0].DecryptData(util.CreateHash(passphrase))
		assert.Equal(t, secret1.Data, plainSecret1.Data)

		// assert both secret is encrypted
		plainSecret2, err := secrets[1].DecryptData(util.CreateHash(passphrase))
		assert.Equal(t, secret2.Data, plainSecret2.Data)
	})
}
