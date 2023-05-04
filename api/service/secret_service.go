package service

import (
	"fmt"

	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/repository"
	"github.com/jinzhu/gorm"
)

type SecretService interface {
	ListSecret(projectID models.ID) ([]*models.Secret, error)
	FindByIDAndProjectID(secretID models.ID, projectID models.ID) (*models.Secret, error)
	Save(secret *models.Secret) (*models.Secret, error)
	Delete(secretID models.ID, projectID models.ID) error
}

func NewSecretService(secretStorage repository.SecretRepository) SecretService {
	return &secretService{
		secretStorage: secretStorage,
	}
}

type secretService struct {
	secretStorage repository.SecretRepository
}

func (ss *secretService) ListSecret(projectID models.ID) ([]*models.Secret, error) {
	return ss.secretStorage.List(projectID)
}

func (ss *secretService) FindByIDAndProjectID(secretID models.ID, projectID models.ID) (*models.Secret, error) {
	secret, err := ss.secretStorage.GetAsPlainText(secretID, projectID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf(
			"error when fetching secret with id: %d, project_id: %d and error: %v",
			secretID, projectID, err)
	}
	return secret, nil
}

func (ss *secretService) Save(secret *models.Secret) (*models.Secret, error) {
	secretFromDB, err := ss.secretStorage.Save(secret)
	if err != nil {
		return nil, fmt.Errorf(
			"error when upsert secret with project_id: %d, name: %v and error: %v",
			secret.ProjectID, secret.Name, err)
	}
	return secretFromDB, nil
}

func (ss *secretService) Delete(secretID models.ID, projectID models.ID) error {
	if err := ss.secretStorage.Delete(secretID, projectID); err != nil {
		return fmt.Errorf(
			"error when deleting secret with id: %d, project_id: %d and error: %v",
			secretID, projectID, err)
	}
	return nil
}
