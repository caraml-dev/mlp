package service

import (
	"github.com/jinzhu/gorm"

	"github.com/gojek/mlp/api/models"
)

type UsersService interface {
	Save(user *models.User) (*models.User, error)
	FindByEmail(userEmail string) (*models.User, error)
}

type usersService struct {
	db *gorm.DB
}

func NewUsersService(db *gorm.DB) UsersService {
	return &usersService{db: db}
}

//Store token here
func (service *usersService) Save(user *models.User) (*models.User, error) {
	//Check that the user is not present in the database
	currentUser, err := service.FindByEmail(user.Email)
	if currentUser == nil {
		if err = service.db.Create(user).Error; err != nil {
			return nil, err
		}
	}
	return service.FindById(user.Id)
}

func (service *usersService) FindById(userId models.Id) (*models.User, error) {
	var user models.User
	if err := service.db.
		Where("users.id = ?", userId).
		First(&user).
		Error; err != nil {

		return nil, err
	}
	return &user, nil
}

func (service *usersService) FindByEmail(userEmail string) (*models.User, error) {
	var user models.User
	if err := service.db.Where("email = ?", userEmail).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
