package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// SecretStorage represents the external secret storage service for storing a secret
type SecretStorage struct {
	// ID is the unique identifier of the secret storage
	ID ID `json:"id"`
	// Name is the name of the secret storage
	Name string `json:"name"`
	// Type is the type of the secret storage
	Type SecretStorageType `json:"type"`
	// Scope of the secret storage, it can be either "global" or "project"
	Scope SecretStorageScope `json:"scope"`
	// ProjectID is the ID of the project that the secret storage belongs to when the scope is "project"
	ProjectID *ID `json:"project_id"`
	// Project is the project that the secret storage belongs to when the scope is "project"
	Project *Project `json:"project,omitempty"`
	// Config is type-specific secret storage configuration
	Config SecretStorageConfig `json:"config"`
	// CreatedUpdated is the timestamp of the creation and last update of the secret storage
	CreatedUpdated
}

type SecretStorageConfig struct {
	// VaultConfig is the configuration of the Vault secret storage. This field is populated when the type is "vault"
	VaultConfig *VaultConfig `json:"vault_config,omitempty"`
}

func (c *SecretStorageConfig) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, c)
}

func (c SecretStorageConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// VaultConfig is the configuration of the Vault secret storage
type VaultConfig struct {
	// Vault URL
	URL string `json:"url"`
	// Role to be used when communicating with Vault
	Role string `json:"role"`
	// Vault kv version
	KVVersion string `json: "kv_version"`
	// MountPath is the path of the secret storage in Vault
	MountPath string `json: "mount_path"`
	// AuthMethod is the authentication method to be used when communicating with Vault
	AuthMethod AuthMethod `json: "auth_method"`
	// GCPAuthType is the GCP authentication type to be used when communicating with Vault, the value can be either "iam" or "gce"
	GCPAuthType *GCPAuthType `json:"gcp_auth_type,omitempty"`
}

// SecretStorageScope is the scope of the secret storage
type SecretStorageScope string

// SecretStorageType is the type of the secret storage
type SecretStorageType string

// AuthMethod is the authentication type to be used when communicating with Vault
type AuthMethod string

// GCPAuthType is the GCP authentication type to be used when communicating with Vault
type GCPAuthType string

const (
	// Secret storage with global scope can be accessed by all projects
	GlobalSecretStorageScope SecretStorageScope = "global"
	// Secret storage with project scope can only be accessed by the project that it belongs to
	ProjectSecretStorageScope SecretStorageScope = "project"

	// InternalSecretStorageType secret storage stores secret in the MLP database
	InternalSecretStorageType SecretStorageType = "internal"
	// VaultSecretStorageType secret storage stores secret in a Vault instance
	VaultSecretStorageType SecretStorageType = "vault"

	// Use gcp authentication method to communicate with Vault https://developer.hashicorp.com/vault/docs/auth/gcp
	GCPAuthMethod AuthMethod = "gcp"
	// Use gce authentication method to communicate with Vault https://developer.hashicorp.com/vault/docs/auth/gcp#gce-login
	GCEGCPAuthType GCPAuthType = "gce"
	// Use iam authentication method to communicate with Vault https://developer.hashicorp.com/vault/docs/auth/gcp#iam-login
	IAMGCPAuthType GCPAuthType = "iam"
)
