package secretstorage

import (
	"fmt"
	"sync"

	"github.com/caraml-dev/mlp/api/models"
)

type Registry struct {
	registry map[models.ID]SecretStorageClient
	lock     sync.RWMutex
}

func NewRegistry(secretStorages []*models.SecretStorage) (*Registry, error) {
	registry := make(map[models.ID]SecretStorageClient)
	for _, ss := range secretStorages {
		if ss.Type == models.InternalSecretStorageType {
			continue
		}

		c, err := NewSecretStorageClient(ss)
		if err != nil {
			return nil, fmt.Errorf("failed to create secret storage vaultClient: %w", err)
		}
		registry[ss.ID] = c
	}

	return &Registry{
		registry: registry,
	}, nil
}

func (r *Registry) Set(secretStorageID models.ID, client SecretStorageClient) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.registry[secretStorageID] = client
}

func (r *Registry) Get(secretStorageID models.ID) (SecretStorageClient, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	sc, ok := r.registry[secretStorageID]
	return sc, ok
}
