package service

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/caraml-dev/mlp/api/models"
	apperror "github.com/caraml-dev/mlp/api/pkg/errors"
	"github.com/caraml-dev/mlp/api/pkg/secretstorage"
	"github.com/caraml-dev/mlp/api/repository"
)

type SecretStorageService interface {
	Create(ss *models.SecretStorage) (*models.SecretStorage, error)
	// FindByID retrieves a secret storage by ID
	FindByID(id models.ID) (*models.SecretStorage, error)
	// List retrieves all secret storages for a project
	List(projectID models.ID) ([]*models.SecretStorage, error)
	// ListAll retrieves all secret storages
	ListAll() ([]*models.SecretStorage, error)
	// Update updates a secret storage
	Update(storage *models.SecretStorage) (*models.SecretStorage, error)
	// Delete deletes a secret storage
	Delete(id models.ID) error
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

func (s *secretStorageService) FindByID(id models.ID) (*models.SecretStorage, error) {
	return s.ssRepository.Get(id)
}

func (s *secretStorageService) List(projectID models.ID) ([]*models.SecretStorage, error) {
	projectSs, err := s.ssRepository.List(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret storage: %w", err)
	}

	globalSs, err := s.ssRepository.ListGlobal()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret storage: %w", err)
	}

	return append(globalSs, projectSs...), nil
}

func (s *secretStorageService) ListAll() ([]*models.SecretStorage, error) {
	return s.ssRepository.ListAll()
}

func (s *secretStorageService) Create(ss *models.SecretStorage) (*models.SecretStorage, error) {
	ss, err := s.ssRepository.Save(ss)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret storage: %w", err)
	}

	// create the client and store it in the registry
	client, err := secretstorage.NewClient(ss)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret storage client: %w", err)
	}
	s.ssClientRegistry.Set(ss.ID, client)

	return ss, nil
}

func (s *secretStorageService) Delete(id models.ID) error {
	ss, err := s.ssRepository.Get(id)
	if err != nil {
		if errors.Is(err, &apperror.NotFoundError{}) {
			return nil
		}
		return fmt.Errorf("failed to retrieve secret storage: %w", err)
	}

	if ss.Type != models.InternalSecretStorageType && ss.Scope == models.GlobalSecretStorageScope {
		return apperror.NewInvalidArgumentErrorf("global secret storage cannot be deleted")
	}

	if ss.Type != models.InternalSecretStorageType {
		client, ok := s.ssClientRegistry.Get(id)
		if !ok {
			return fmt.Errorf("secret storage client not found")
		}

		err = client.DeleteAll(ss.Project.Name)
		if err != nil {
			return fmt.Errorf("failed to delete secrets in secret storage: %w", err)
		}
	}

	return s.ssRepository.Delete(id)
}

func (s *secretStorageService) Update(ss *models.SecretStorage) (*models.SecretStorage, error) {
	existingSs, err := s.ssRepository.Get(ss.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret storage: %w", err)
	}

	if existingSs.Type != ss.Type || !reflect.DeepEqual(existingSs.Config, ss.Config) {
		return s.migrateSecretStorage(existingSs, ss)
	}

	return s.ssRepository.Save(ss)
}

func (s *secretStorageService) migrateSecretStorage(oldSs *models.SecretStorage,
	newSs *models.SecretStorage) (*models.SecretStorage, error) {

	client, ok := s.ssClientRegistry.Get(oldSs.ID)
	if !ok {
		return nil, fmt.Errorf("secret storage client not found")
	}

	allSecrets, err := client.List(oldSs.Project.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets in secret storage: %w", err)
	}

	newClient, err := secretstorage.NewClient(newSs)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret storage client: %w", err)
	}

	err = newClient.SetAll(allSecrets, newSs.Project.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to set secrets in secret storage: %w", err)
	}

	err = client.DeleteAll(oldSs.Project.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to delete secrets in secret storage: %w", err)
	}

	return s.ssRepository.Save(newSs)
}
