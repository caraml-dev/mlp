package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/gojek/mlp/api/models"
	"github.com/gojek/mlp/api/util"
)

type SecretStorage interface {
	// List list all secret within the given project ID.
	// The secrets returned is in encrypted form
	List(projectID models.ID) ([]*models.Secret, error)

	// Save encrypt and save a plain text secret.
	// The secrets returned is in encrypted form
	Save(secret *models.Secret) (*models.Secret, error)

	// GetAsPlainText return a secret in plain text
	GetAsPlainText(id models.ID, projectID models.ID) (*models.Secret, error)

	// GetByNameAsPlainText return a secret in plain text
	GetByNameAsPlainText(name string, projectID models.ID) (*models.Secret, error)

	// Delete delete secret
	Delete(id models.ID, projectID models.ID) error
}

type secretStorage struct {
	db            *gorm.DB
	encryptionKey string
}

func NewSecretStorage(db *gorm.DB, passphrase string) SecretStorage {
	return &secretStorage{db: db,
		encryptionKey: util.CreateHash(passphrase),
	}
}

func (ss *secretStorage) List(projectID models.ID) (secrets []*models.Secret, err error) {
	err = ss.db.Where("project_id = ?", projectID).Find(&secrets).Error
	return
}

func (ss *secretStorage) Save(secret *models.Secret) (*models.Secret, error) {
	encSecret, err := secret.EncryptData(ss.encryptionKey)

	if err != nil {
		return nil, fmt.Errorf(
			"error when decrypt secret data with project_id: %d, name: %s and error: %v",
			secret.ProjectID, secret.Name, err)
	}

	if err := ss.db.Save(encSecret).Error; err != nil {
		return nil, err
	}
	return encSecret, nil
}

func (ss *secretStorage) Delete(id models.ID, projectID models.ID) error {
	return ss.db.Where("id = ? AND project_id = ?", id, projectID).Delete(models.Secret{}).Error
}

func (ss *secretStorage) GetAsPlainText(id models.ID, projectID models.ID) (*models.Secret, error) {
	var secret models.Secret
	if err := ss.db.Where("id = ? AND project_id = ?", id, projectID).First(&secret).Error; err != nil {
		return nil, err
	}
	decSecret, err := secret.DecryptData(ss.encryptionKey)

	if err != nil {
		return nil, fmt.Errorf("error when decrypt secret data with id: %s with error: %v", id, err)
	}
	return decSecret, nil
}

func (ss *secretStorage) GetByNameAsPlainText(name string, projectID models.ID) (*models.Secret, error) {
	var secret models.Secret
	if err := ss.db.Where("project_id = ? AND name = ?", projectID, name).First(&secret).Error; err != nil {
		return nil, err
	}

	decSecret, err := secret.DecryptData(ss.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error when decrypt secret data with name: %s with error: %v", name, err)
	}
	return decSecret, nil
}
