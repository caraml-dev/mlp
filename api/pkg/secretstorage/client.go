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
	// Delete deletes a CaraML secret from the secret storage
	Delete(name string, project string) error
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
