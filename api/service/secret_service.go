package service

import (
	"fmt"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/secretstorage"
	"github.com/caraml-dev/mlp/api/repository"
	"github.com/caraml-dev/mlp/api/util"
)

// SecretService is the interface that provides secret related methods.
type SecretService interface {
	// Create creates a secret in the storage and returns the created secret.
	Create(secret *models.Secret) (*models.Secret, error)
	// Update updates a secret in the storage and returns the updated secret.
	Update(secret *models.Secret) (*models.Secret, error)
	// List lists all secrets of a project given its projectID
	List(projectID models.ID) ([]*models.Secret, error)
	// Delete deletes a secret given its secretID
	Delete(secretID models.ID) error
}

func NewSecretService(secretRepository repository.SecretRepository,
	storageRepository repository.SecretStorageRepository,
	projectRepository repository.ProjectRepository,
	storageClientRegistry *secretstorage.Registry,
	defaultSecretStorage *models.SecretStorage,
) SecretService {
	return &secretService{
		secretRepository:  secretRepository,
		storageRepository: storageRepository,
		projectRepository: projectRepository,

		storageClientRegistry: storageClientRegistry,
		defaultSecretStorage:  defaultSecretStorage,
	}
}

type secretService struct {
	secretRepository  repository.SecretRepository
	storageRepository repository.SecretStorageRepository
	projectRepository repository.ProjectRepository

	storageClientRegistry *secretstorage.Registry
	defaultSecretStorage  *models.SecretStorage
}

// Create creates a secret in the storage and returns the created secret.
func (ss *secretService) Create(secret *models.Secret) (*models.Secret, error) {
	project, err := ss.projectRepository.Get(secret.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("error when fetching project with id: %d, error: %w", secret.ProjectID, err)
	}

	// get secret storage, use default if users don't specify
	secretStorage := ss.defaultSecretStorage
	if secret.SecretStorageID != nil {
		secretStorage, err = ss.storageRepository.Get(*secret.SecretStorageID)
		if err != nil {
			return nil, fmt.Errorf("error when fetching secret storage with id: %d, error: %w", *secret.SecretStorageID, err)
		}
	} else {
		secret.SecretStorageID = &ss.defaultSecretStorage.ID
	}

	// for internal secret we can simply store to DB
	if secretStorage.Type == models.InternalSecretStorageType {
		// create secret in database, including the data
		return ss.secretRepository.Save(secret)
	}

	// Get the corresponding secret storage client
	ssClient, ok := ss.storageClientRegistry.Get(secretStorage.ID)
	if !ok {
		return nil, fmt.Errorf("secret storage client with id %d is not found", secretStorage.ID)
	}

	// Update secret data in the corresponding secret storage
	err = ssClient.Set(secret.Name, secret.Data, project.Name)
	if err != nil {
		return nil, fmt.Errorf("error when creating secret in secret storage with id: %d, error: %w", *secret.SecretStorageID, err)
	}

	return ss.saveExternalSecret(secret)
}

// List lists all secrets of a project given its projectID
func (ss *secretService) List(projectID models.ID) ([]*models.Secret, error) {
	project, err := ss.projectRepository.Get(projectID)
	if err != nil {
		return nil, fmt.Errorf("error when fetching project with id: %d, error: %w", projectID, err)
	}

	secrets, err := ss.secretRepository.List(projectID)
	if err != nil {
		return nil, fmt.Errorf("error when fetching secrets with project_id: %d, error: %w", projectID, err)
	}

	// group secrets by storage ID, skip internal storage since its 'data' is directly available in the secret object
	// grouping it by storage ID so that we don't make multiple calls to the secret storage
	secretsByStorageID := make(map[models.ID][]*models.Secret)
	for _, secret := range secrets {
		if secret.SecretStorage.Type == models.InternalSecretStorageType {
			continue
		}

		secretsByStorageID[*secret.SecretStorageID] = append(secretsByStorageID[*secret.SecretStorageID], secret)
	}

	// fetch secrets from secret storage
	secretKVs := make(map[string]string)
	for storageID := range secretsByStorageID {
		secretStorageClient, ok := ss.storageClientRegistry.Get(storageID)
		if !ok {
			return nil, fmt.Errorf("secret storage client with id %d is not found", storageID)
		}

		temp, err := secretStorageClient.List(project.Name)
		if err != nil {
			return nil, fmt.Errorf("error when fetching secrets from secret storage with id: %d, error: %w", storageID, err)
		}

		secretKVs = util.JoinMaps(secretKVs, temp)
	}

	// populate 'data' field of secrets
	for _, secret := range secrets {
		if secret.SecretStorage.Type == models.InternalSecretStorageType {
			continue
		}

		secret.Data = secretKVs[secret.Name]
	}

	return ss.secretRepository.List(projectID)
}

func (ss *secretService) Update(secret *models.Secret) (*models.Secret, error) {
	project, err := ss.projectRepository.Get(secret.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("error when fetching project with id: %d, error: %w", secret.ProjectID, err)
	}

	existingSecret, err := ss.secretRepository.Get(secret.ID)
	if err != nil {
		return nil, fmt.Errorf("error when fetching secret with id: %d, project_id: %d, error: %w", secret.ID, secret.ProjectID, err)
	}

	// secret storage id is changed, migrate the secret to the new storage
	if *secret.SecretStorageID != existingSecret.SecretStorage.ID {
		return ss.migrateSecret(existingSecret, secret, project)
	}

	secretStorage, err := ss.storageRepository.Get(*secret.SecretStorageID)
	if err != nil {
		return nil, fmt.Errorf("error when fetching secret storage with id: %d, error: %w", *secret.SecretStorageID, err)
	}

	if secretStorage.Type == models.InternalSecretStorageType {
		// create secret in database, including the data
		return ss.secretRepository.Save(secret)
	}

	// Get the corresponding secret storage client
	ssClient, ok := ss.storageClientRegistry.Get(secretStorage.ID)
	if !ok {
		return nil, fmt.Errorf("secret storage client with id %d is not found", secretStorage.ID)
	}

	// Update secret data in the corresponding secret storage
	err = ssClient.Set(project.Name, secret.Data, secret.Name)
	if err != nil {
		return nil, fmt.Errorf("error when creating secret in secret storage with id: %d, error: %w", *secret.SecretStorageID, err)
	}

	return ss.saveExternalSecret(secret)
}

func (ss *secretService) Delete(secretID models.ID) error {
	if err := ss.secretRepository.Delete(secretID); err != nil {
		return fmt.Errorf(
			"error when deleting secret with id: %d, error: %w",
			secretID, err)
	}
	return nil
}

// migrateSecret migrate secret from one secret storage to another
func (ss *secretService) migrateSecret(oldSecret *models.Secret, newSecret *models.Secret, project *models.Project) (*models.Secret, error) {
	newSecretStorage, err := ss.storageRepository.Get(*newSecret.SecretStorageID)
	if err != nil {
		return nil, fmt.Errorf("error when fetching secret storage with id: %d, error: %w", *newSecret.SecretStorageID, err)
	}

	// disallow migrating to internal secret storage
	if newSecretStorage.Type == models.InternalSecretStorageType {
		return nil, fmt.Errorf("cannot migrate secret to internal secret storage")
	}

	oldSecretStorage, err := ss.storageRepository.Get(*oldSecret.SecretStorageID)
	if err != nil {
		return nil, fmt.Errorf("error when fetching secret storage with id: %d, error: %w", *oldSecret.SecretStorageID, err)
	}

	// for internal secret type, "oldSecret.Data" already stores the secret value
	// for non-internal secret type we'll have to fetch it from corresponding secret storage
	if oldSecretStorage.Type != models.InternalSecretStorageType {
		oldSsClient, ok := ss.storageClientRegistry.Get(oldSecretStorage.ID)
		if !ok {
			return nil, fmt.Errorf("secret storage client with id %d is not found", oldSecretStorage.ID)
		}

		secretValue, err := oldSsClient.Get(oldSecret.Name, project.Name)
		if err != nil {
			return nil, fmt.Errorf("error when fetching secret from secret storage with id: %d, error: %w", *oldSecret.SecretStorageID, err)
		}

		oldSecret.Data = secretValue

		// delete secret from old storage
		err = oldSsClient.Delete(oldSecret.Name, project.Name)
		if err != nil {
			return nil, fmt.Errorf("error when deleting secret from secret storage with id: %d, error: %w", *oldSecret.SecretStorageID, err)
		}
	}

	if newSecret.Data == "" {
		newSecret.Data = oldSecret.Data
	}

	newSsClient, ok := ss.storageClientRegistry.Get(newSecretStorage.ID)
	if !ok {
		return nil, fmt.Errorf("secret storage client with id %d is not found", newSecretStorage.ID)
	}

	err = newSsClient.Set(newSecret.Name, newSecret.Data, project.Name)
	if err != nil {
		return nil, fmt.Errorf("error when creating secret in secret storage with id: %d, error: %w", *newSecret.SecretStorageID, err)
	}

	return ss.saveExternalSecret(newSecret)
}

func (ss *secretService) saveExternalSecret(secret *models.Secret) (*models.Secret, error) {
	secretData := secret.Data
	secret.Data = "" // don't store secret data in DB for external secret
	secret, err := ss.secretRepository.Save(secret)
	secret.Data = secretData
	return secret, err
}
