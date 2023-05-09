package repository

import (
	"github.com/caraml-dev/mlp/api/models"
	"github.com/jinzhu/gorm"
)

// SecretStorageRepository is an interface for interacting with "secret_storages" table in DB
type SecretStorageRepository interface {
	// Get returns a Secret Storage with a name within a project
	Get(name string, project string) (*models.SecretStorage, error)
	// List lists all Secret Storage within a project
	List(project string) ([]*models.SecretStorage, error)
	// GetGlobal return a global Secret Storage with a name
	GetGlobal(name string) (*models.SecretStorage, error)
	// ListGlobal lists all global Secret Storage
	ListGlobal() ([]*models.SecretStorage, error)
	// Save creates or updates a Secret Storage
	Save(secretStorage *models.SecretStorage) (*models.SecretStorage, error)
	// Delete deletes a Secret Storage
	Delete(secretStorage *models.SecretStorage) error
}

type secretStorageRepository struct {
	db *gorm.DB
}

// NewSecretStorageRepository creates a new Secret Storage Repository
func NewSecretStorageRepository(db *gorm.DB) SecretStorageRepository {
	return &secretStorageRepository{
		db: db,
	}
}

// Get returns a Secret Storage with a name within a project
func (r *secretStorageRepository) Get(name string, project string) (*models.SecretStorage, error) {
	var ss models.SecretStorage

	err := r.db.Table("secret_storages").
		Preload("Project").
		Joins("JOIN projects ON secret_storages.project_id = projects.id").
		Where("projects.name = ? AND secret_storages.name = ?", project, name).
		First(&ss).Error

	return &ss, err
}

// List lists all Secret Storage within a project
func (r *secretStorageRepository) List(project string) ([]*models.SecretStorage, error) {
	var ss []*models.SecretStorage

	err := r.db.Table("secret_storages").
		Preload("Project").
		Joins("JOIN projects ON secret_storages.project_id = projects.id").
		Where("projects.name = ?", project).
		Find(&ss).Error

	return ss, err
}

// GetGlobal return a global Secret Storage with a name
func (r *secretStorageRepository) GetGlobal(name string) (*models.SecretStorage, error) {
	var ss models.SecretStorage

	err := r.db.Where("name = ? AND scope = ?", name, models.GlobalSecretStorageScope).Find(&ss).Error
	return &ss, err
}

// ListGlobal lists all global Secret Storage
func (r *secretStorageRepository) ListGlobal() ([]*models.SecretStorage, error) {
	var ss []*models.SecretStorage

	err := r.db.Where("scope = ?", models.GlobalSecretStorageScope).Find(&ss).Error
	return ss, err
}

// Save creates or updates a Secret Storage
func (r *secretStorageRepository) Save(secretStorage *models.SecretStorage) (*models.SecretStorage, error) {
	err := r.db.Save(secretStorage).Error
	return secretStorage, err
}

// Delete deletes a Secret Storage
func (r *secretStorageRepository) Delete(secretStorage *models.SecretStorage) error {
	return r.db.Delete(secretStorage).Error
}