package service

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gojek/mlp/api/models"
	"github.com/gojek/mlp/api/storage/mocks"
)

func TestFindByIdAndProjectId(t *testing.T) {
	testCases := []struct {
		desc             string
		secretFromDB     *models.Secret
		errorFetchFromDb error
		expectedSecret   *models.Secret
		expectedError    string
	}{
		{
			desc: "Should success",
			secretFromDB: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "qzSQ+pZ9Qu7+SpTQCuZB2AgdtH3cuMR0eWbH/yvlqrI=",
			},
			expectedSecret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "qzSQ+pZ9Qu7+SpTQCuZB2AgdtH3cuMR0eWbH/yvlqrI=",
			},
		},
		{
			desc: "Should return nil and no error if record not found",
			secretFromDB: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "qzSQ+pZ9Qu7+SpTQCuZB2AgdtH3cuMR0eWbH/yvlqrI=",
			},
			errorFetchFromDb: gorm.Errors{gorm.ErrRecordNotFound},
			expectedSecret:   nil,
		},
		{
			desc: "Should return error if something going wrong when fetching db",
			secretFromDB: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "qzSQ+pZ9Qu7+SpTQCuZB2AgdtH3cuMR0eWbH/yvlqrI=",
			},
			errorFetchFromDb: fmt.Errorf("db is down"),
			expectedError:    "error when fetching secret with id: 1, project_id: 1 and error: db is down",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			secretStorage := &mocks.SecretStorage{}
			secretStorage.On("GetAsPlainText", models.Id(1), models.Id(1)).Return(tC.secretFromDB, tC.errorFetchFromDb)
			secretService := NewSecretService(secretStorage)
			result, err := secretService.FindByIdAndProjectId(1, 1)
			if tC.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, tC.expectedSecret, result)
			} else {
				assert.EqualError(t, err, tC.expectedError)
			}
		})
	}
}

func TestSave(t *testing.T) {
	testCases := []struct {
		desc          string
		secret        *models.Secret
		errorFromDB   error
		expectedError string
	}{
		{
			desc: "Should success",
			secret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "plainData",
			},
		},
		{
			desc: "Should raise error when failed save to db",
			secret: &models.Secret{
				Id:        models.Id(1),
				ProjectId: models.Id(1),
				Name:      "name",
				Data:      "plainData",
			},
			errorFromDB:   fmt.Errorf("db is down"),
			expectedError: "error when upsert secret with project_id: 1, name: name and error: db is down",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			secretStorage := &mocks.SecretStorage{}
			secretStorage.On("Save", tC.secret).Return(tC.secret, tC.errorFromDB)
			secretService := NewSecretService(secretStorage)
			result, err := secretService.Save(tC.secret)
			if tC.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, tC.secret.Id, result.Id)
				assert.Equal(t, tC.secret.Name, result.Name)
				assert.Equal(t, tC.secret.ProjectId, result.ProjectId)
				require.NoError(t, err)
			} else {
				assert.EqualError(t, err, tC.expectedError)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		desc          string
		secretId      models.Id
		projectId     models.Id
		errorFromDB   error
		expectedError string
	}{
		{
			desc:      "Should success",
			secretId:  models.Id(1),
			projectId: models.Id(1),
		},
		{
			desc:          "Should success",
			secretId:      models.Id(1),
			projectId:     models.Id(1),
			errorFromDB:   fmt.Errorf("db is down"),
			expectedError: "error when deleting secret with id: 1, project_id: 1 and error: db is down",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			secretStorage := &mocks.SecretStorage{}
			secretStorage.On("Delete", tC.secretId, tC.projectId).Return(tC.errorFromDB)
			secretService := &secretService{
				secretStorage: secretStorage,
			}
			err := secretService.Delete(tC.secretId, tC.projectId)
			if tC.expectedError == "" {
				require.NoError(t, err)
			} else {
				assert.EqualError(t, err, tC.expectedError)
			}
		})
	}
}

func TestList(t *testing.T) {
	projectId := models.Id(1)
	secrets := []*models.Secret{
		{
			Id:        models.Id(1),
			ProjectId: projectId,
			Name:      "name1",
			Data:      "plainData",
		},
		{
			Id:        models.Id(2),
			ProjectId: projectId,
			Name:      "name2",
			Data:      "plainData",
		},
	}

	secretStorage := &mocks.SecretStorage{}
	secretStorage.On("List", projectId).Return(secrets, nil)
	secretService := &secretService{
		secretStorage: secretStorage,
	}
	actual, err := secretService.ListSecret(projectId)
	assert.NoError(t, err)
	assert.Equal(t, secrets, actual)
}
