package service

import (
	"github.com/jinzhu/gorm"

	"github.com/gojek/mlp/models"
)

type AccountService interface {
	FindByUserId(userId models.Id) (*models.Account, error)
	Save(token *models.Account) (*models.Account, error)
}

type accountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) AccountService {
	return &accountService{db: db}
}

//Store token here
func (service *accountService) Save(account *models.Account) (*models.Account, error) {
	if err := service.db.Create(account).Error; err != nil {
		return nil, err
	}
	return service.FindByUserId(account.UserId)
}

func (service *accountService) FindByUserId(userId models.Id) (*models.Account, error) {
	var userAccount models.Account
	if err := service.db.Where("user_id = ?", userId).First(&userAccount).Error; err != nil {
		return nil, err
	}
	return &userAccount, nil
}
