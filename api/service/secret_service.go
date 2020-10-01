package service

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/gojek/mlp/api/models"
	"github.com/gojek/mlp/api/storage"
)

type SecretService interface {
	ListSecret(projectId models.Id) ([]*models.Secret, error)
	FindByIdAndProjectId(secretId models.Id, projectId models.Id) (*models.Secret, error)
	Save(secret *models.Secret) (*models.Secret, error)
	Delete(secretId models.Id, projectId models.Id) error
}

func NewSecretService(secretStorage storage.SecretStorage) SecretService {
	return &secretService{
		secretStorage: secretStorage,
	}
}

type secretService struct {
	secretStorage storage.SecretStorage
}

func (ss *secretService) ListSecret(projectId models.Id) ([]*models.Secret, error) {
	return ss.secretStorage.List(projectId)
}

func (ss *secretService) FindByIdAndProjectId(secretId models.Id, projectId models.Id) (*models.Secret, error) {
	secret, err := ss.secretStorage.GetAsPlainText(secretId, projectId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error when fetching secret with id: %d, project_id: %d and error: %v", secretId, projectId, err)
	}
	return secret, nil
}

func (ss *secretService) Save(secret *models.Secret) (*models.Secret, error) {
	secretFromDB, err := ss.secretStorage.Save(secret)
	if err != nil {
		return nil, fmt.Errorf("error when upsert secret with project_id: %d, name: %v and error: %v", secret.ProjectId, secret.Name, err)
	}
	return secretFromDB, nil
}

func (ss *secretService) Delete(secretId models.Id, projectId models.Id) error {
	if err := ss.secretStorage.Delete(secretId, projectId); err != nil {
		return fmt.Errorf("error when deleting secret with id: %d, project_id: %d and error: %v", secretId, projectId, err)
	}
	return nil
}
