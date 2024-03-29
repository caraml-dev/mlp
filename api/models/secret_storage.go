package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
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
	ProjectID *ID `json:"project_id,omitempty"`
	// Project is the project that the secret storage belongs to when the scope is "project"
	Project *Project `json:"-"`
	// Config is type-specific secret storage configuration
	Config SecretStorageConfig `json:"config,omitempty"`
	// CreatedUpdated is the timestamp of the creation and last update of the secret storage
	CreatedUpdated
}

func (s *SecretStorage) ValidateForCreation() error {
	return s.validate(false)
}

func (s *SecretStorage) ValidateForMutation() error {
	return s.validate(true)
}

func (s *SecretStorage) MergeValue(other *SecretStorage) error {
	return copier.CopyWithOption(s, other, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

func (s *SecretStorage) validate(checkID bool) error {
	if checkID && s.ID <= 0 {
		return fmt.Errorf("invalid secret storage ID: %d", s.ID)
	}

	maxNameChar := 64
	if s.Name == "" || len(s.Name) > maxNameChar {
		return fmt.Errorf("invalid secret storage name: %s", s.Name)
	}

	if s.Type == "" || s.Type == InternalSecretStorageType {
		return fmt.Errorf("invalid secret storage type: %s", s.Type)
	}

	if s.Scope != ProjectSecretStorageScope {
		return fmt.Errorf("invalid secret storage scope: %s", s.Scope)
	}

	if s.ProjectID == nil {
		return fmt.Errorf("invalid secret storage project ID: %d", s.ProjectID)
	}

	return nil
}

type SecretStorageConfig struct {
	// VaultConfig is the configuration of the Vault secret storage.
	// This field is populated when the type is "vault"
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
	// MountPath is the path of the secret storage in Vault
	MountPath string `json:"mount_path"`
	// PathPrefix is the prefix of the path of the secret in Vault
	PathPrefix string `json:"path_prefix"`
	// AuthMethod is the authentication method to be used when communicating with Vault
	AuthMethod AuthMethod `json:"auth_method"`
	// GCPAuthType is the GCP authentication type to be used when communicating with Vault.
	// The value can be either "iam" or "gce"
	GCPAuthType GCPAuthType `json:"gcp_auth_type,omitempty"`
	// Token is the token to be used when communicating with Vault
	// This field is only used when the auth method is "token"
	// Only use this method when Vault is running in dev mode
	Token string `json:"token,omitempty"`
	// ServiceAccountEmail is the service account email to be used when communicating with Vault
	// This field is only used when the AuthMethod is "gcp" and GCPAuthType is "iam"
	ServiceAccountEmail string `json:"service_account_email"`
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

	// Use gcp authentication method to communicate with Vault
	// https://developer.hashicorp.com/vault/docs/auth/gcp
	GCPAuthMethod AuthMethod = "gcp"
	// Use gce authentication method to communicate with Vault
	// https://developer.hashicorp.com/vault/docs/auth/gcp#gce-login
	GCEGCPAuthType GCPAuthType = "gce"
	// Use iam authentication method to communicate with Vault
	// https://developer.hashicorp.com/vault/docs/auth/gcp#iam-login
	IAMGCPAuthType GCPAuthType = "iam"

	// Use token authentication method to communicate with Vault
	// Only use this method when Vault is running in dev mode
	TokenAuthMethod AuthMethod = "token"
)
