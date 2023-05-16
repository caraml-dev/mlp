package secretstorage

import (
	"fmt"

	"github.com/caraml-dev/mlp/api/models"
)

type Client interface {
	// Get retrieves a CaraML secret from the secret storage
	Get(name string, project string) (string, error)
	// Set creates or updates a CaraML secret in the secret storage
	Set(name string, secretValue string, project string) error
	// List lists all CaraML secrets in the secret storage
	List(project string) (map[string]string, error)
	// SetAll creates or updates all CaraML secrets of a project in the secret storage
	SetAll(secrets map[string]string, project string) error
	// Delete deletes a CaraML secret from the secret storage
	Delete(name string, project string) error
	// DeleteAll deletes all CaraML secrets from the secret storage
	DeleteAll(project string) error
}

// NewClient creates a new secret storage client
func NewClient(ss *models.SecretStorage) (Client, error) {
	switch ss.Type {
	case models.VaultSecretStorageType:
		return NewVaultSecretStorageClient(ss)
	default:
		return nil, fmt.Errorf("unsupported secret storage type %s", ss.Type)
	}
}
