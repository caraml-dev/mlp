package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/caraml-dev/mlp/api/models"
)

type ProjectRepository interface {
	ListProjects(name string) ([]*models.Project, error)
	Get(projectID models.ID) (*models.Project, error)
	GetByName(projectName string) (*models.Project, error)
	Save(project *models.Project) (*models.Project, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (storage *projectRepository) ListProjects(name string) (projects []*models.Project, err error) {
	err = storage.db.Where("name LIKE ?", name+"%").Find(&projects).Error
	return
}

func (storage *projectRepository) Get(projectID models.ID) (*models.Project, error) {
	var project models.Project
	if err := storage.db.Where("id = ?", projectID).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (storage *projectRepository) GetByName(projectName string) (*models.Project, error) {
	var project models.Project
	if err := storage.db.Where("name = ?", projectName).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (storage *projectRepository) Save(project *models.Project) (*models.Project, error) {
	if err := storage.db.Save(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}
