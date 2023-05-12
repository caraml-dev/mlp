package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/caraml-dev/mlp/api/models"
)

type SecretRepository interface {
	// Get return a secret given the secret id
	Get(id models.ID) (*models.Secret, error)
	// List lists all secret within the given project ID.
	List(projectID models.ID) ([]*models.Secret, error)
	// Save create or update a secret.
	Save(secret *models.Secret) (*models.Secret, error)
	// Delete delete secret given the secret id
	Delete(id models.ID) error
}

type secretRepository struct {
	db *gorm.DB
}

func NewSecretRepository(db *gorm.DB) SecretRepository {
	return &secretRepository{
		db: db,
	}
}

// List lists all secret within the given project ID.
func (ss *secretRepository) List(projectID models.ID) (secrets []*models.Secret, err error) {
	err = ss.db.Where("project_id = ?", projectID).Find(&secrets).Error
	return
}

// Save create or update a secret.
func (ss *secretRepository) Save(secret *models.Secret) (*models.Secret, error) {
	if err := ss.db.Save(secret).Error; err != nil {
		return nil, err
	}
	return secret, nil
}

// Delete delete secret given the secret id
func (ss *secretRepository) Delete(id models.ID) error {
	return ss.db.Where("id = ?", id).Delete(models.Secret{}).Error
}

// Get return a secret given the secret id
func (ss *secretRepository) Get(id models.ID) (*models.Secret, error) {
	var secret models.Secret
	if err := ss.db.Where("id = ?", id).First(&secret).Error; err != nil {
		return nil, err
	}

	return &secret, nil
}
