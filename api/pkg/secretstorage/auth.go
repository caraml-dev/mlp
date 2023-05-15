package secretstorage

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	gcpauth "github.com/hashicorp/vault/api/auth/gcp"

	"github.com/caraml-dev/mlp/api/models"
)

// authHelper is an interface for authentication against Vault
type authHelper interface {
	// login authenticates against Vault and returns a client token
	login(client *vault.Client) (*vault.Secret, error)
}

// gcpAuthHelper is an implementation of authHelper for GCP auth
type gcpAuthHelper struct {
	role        string
	loginOption gcpauth.LoginOption
}

// newGcpAuthHelper creates a new GCP auth helper
func newGcpAuthHelper(vaultConfig *models.VaultConfig) authHelper {
	var loginOption gcpauth.LoginOption
	switch vaultConfig.GCPAuthType {
	case models.GCEGCPAuthType:
		loginOption = gcpauth.WithGCEAuth()
	case models.IAMGCPAuthType:
		loginOption = gcpauth.WithIAMAuth(vaultConfig.ServiceAccountEmail)
	}

	return &gcpAuthHelper{
		role:        vaultConfig.Role,
		loginOption: loginOption,
	}
}

// login authenticates against Vault and returns a client token
func (g *gcpAuthHelper) login(client *vault.Client) (*vault.Secret, error) {
	gcpAuth, err := gcpauth.NewGCPAuth(
		g.role,
		g.loginOption,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to initialize GCP auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), gcpAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to GCP auth method: %w", err)
	}

	if authInfo == nil {
		return nil, fmt.Errorf("login response did not return client token")
	}

	return authInfo, nil
}
