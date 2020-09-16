package service

import (
	"github.com/jinzhu/gorm"

	"github.com/gojek/mlp/models"
)

type ApplicationService interface {
	List() ([]*models.Application, error)
}

func NewApplicationService(db *gorm.DB) (ApplicationService, error) {
	return &applicationService{
		DB: db,
	}, nil
}

type applicationService struct {
	*gorm.DB
}

func (service *applicationService) List() (apps []*models.Application, err error) {
	err = service.Where("is_disabled = FALSE").Find(&apps).Error
	return
}
