package service

import (
	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/secretstorage"
	"github.com/caraml-dev/mlp/api/repository"
)

type SecretStorageService interface {
	FindByID(id models.ID) (*models.SecretStorage, error)
	ListAll() ([]*models.SecretStorage, error)
}

type secretStorageService struct {
	ssRepository     repository.SecretStorageRepository
	ssClientRegistry *secretstorage.Registry
}

func NewSecretStorageService(ssRepository repository.SecretStorageRepository,
	ssClientRegistry *secretstorage.Registry) SecretStorageService {
	return &secretStorageService{
		ssRepository:     ssRepository,
		ssClientRegistry: ssClientRegistry,
	}
}

func (s secretStorageService) FindByID(id models.ID) (*models.SecretStorage, error) {
	return s.ssRepository.Get(id)
}

func (s secretStorageService) ListAll() ([]*models.SecretStorage, error) {
	return s.ssRepository.ListAll()
}
