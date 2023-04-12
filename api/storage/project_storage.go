package storage

import (
	"github.com/jinzhu/gorm"

	"github.com/caraml-dev/mlp/api/models"
)

type ProjectStorage interface {
	ListProjects(name string) ([]*models.Project, error)
	Get(projectID models.ID) (*models.Project, error)
	GetByName(projectName string) (*models.Project, error)
	Save(project *models.Project) (*models.Project, error)
}

type projectStorage struct {
	db *gorm.DB
}

func NewProjectStorage(db *gorm.DB) ProjectStorage {
	return &projectStorage{db: db}
}

func (storage *projectStorage) ListProjects(name string) (projects []*models.Project, err error) {
	err = storage.db.Where("name LIKE ?", name+"%").Find(&projects).Error
	return
}

func (storage *projectStorage) Get(projectID models.ID) (*models.Project, error) {
	var project models.Project
	if err := storage.db.Where("id = ?", projectID).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (storage *projectStorage) GetByName(projectName string) (*models.Project, error) {
	var project models.Project
	if err := storage.db.Where("name = ?", projectName).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (storage *projectStorage) Save(project *models.Project) (*models.Project, error) {
	if err := storage.db.Save(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}
