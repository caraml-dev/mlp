package secretstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/caraml-dev/mlp/api/models"
	mlperror "github.com/caraml-dev/mlp/api/pkg/errors"
	vault "github.com/hashicorp/vault/api"
)

type vaultSecretStorageClient struct {
	secretPathTemplate *template.Template
	vaultClient        *vault.Client
	vaultConfig        *models.VaultConfig
}

func NewVaultSecretStorageClient(ss *models.SecretStorage) (SecretStorageClient, error) {
	tmpl, err := template.New("secret_path").Parse(ss.Config.VaultConfig.PathPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to parse secret path template: %w", err)
	}

	// create vault client
	vc := vault.DefaultConfig()
	vc.Address = ss.Config.VaultConfig.URL
	vaultClient, err := vault.NewClient(vc)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	cli := &vaultSecretStorageClient{
		secretPathTemplate: tmpl,
		vaultConfig:        ss.Config.VaultConfig,
		vaultClient:        vaultClient,
	}

	// TODO: call auth method based on ss.Config.VaultConfig.AuthMethod
	// TODO: refresh token periodically
	err = cli.login()
	if err != nil {
		return nil, fmt.Errorf("failed to login to vault: %w", err)
	}

	return cli, nil
}

// Get retrieves a CaraML secret from Vault
// secret is stored as subkey of secretPath
func (v *vaultSecretStorageClient) Get(name string, project string) (string, error) {
	secretPath, err := v.secretPath(project)
	if err != nil {
		return "", err
	}

	secret, err := v.vaultClient.KVv2(v.vaultConfig.MountPath).Get(context.Background(), secretPath)
	if err != nil {
		if errors.Is(err, vault.ErrSecretNotFound) {
			return "", mlperror.NewNotFoundErrorf("secret %s not found in project %s", name, project)
		}

		return "", err
	}

	secretData, ok := secret.Data[name]
	if !ok || secretData == nil {
		return "", mlperror.NewNotFoundErrorf("secret %s not found in project %s", name, project)
	}

	return secretData.(string), nil
}

func (v *vaultSecretStorageClient) Set(name string, secretValue string, project string) error {
	secretPath, err := v.secretPath(project)
	if err != nil {
		return err
	}

	var existingSecrets map[string]interface{}
	secret, err := v.vaultClient.KVv2(v.vaultConfig.MountPath).Get(context.Background(), secretPath)
	if err != nil {
		if !errors.Is(err, vault.ErrSecretNotFound) {
			return err
		}
		existingSecrets = make(map[string]interface{})
	} else {
		existingSecrets = secret.Data
	}

	existingSecrets[name] = secretValue
	_, err = v.vaultClient.KVv2(v.vaultConfig.MountPath).Put(context.Background(), secretPath, existingSecrets)

	return err
}

func (v *vaultSecretStorageClient) List(project string) (map[string]string, error) {
	secretPath, err := v.secretPath(project)
	if err != nil {
		return nil, err
	}

	secretMap := make(map[string]string)
	secret, err := v.vaultClient.KVv2(v.vaultConfig.MountPath).Get(context.Background(), secretPath)
	if err != nil {
		if errors.Is(err, vault.ErrSecretNotFound) {
			return secretMap, nil
		}
		return nil, err
	}

	for k, v := range secret.Data {
		secretMap[k] = v.(string)
	}

	return secretMap, nil
}

func (v *vaultSecretStorageClient) Delete(name string, project string) error {
	secretPath, err := v.secretPath(project)
	if err != nil {
		return err
	}

	secret, err := v.vaultClient.KVv2(v.vaultConfig.MountPath).Get(context.Background(), secretPath)
	if err != nil {
		if errors.Is(err, vault.ErrSecretNotFound) {
			return nil
		}
		return err
	}

	secret.Data[name] = nil

	_, err = v.vaultClient.KVv2(v.vaultConfig.MountPath).Put(context.Background(), secretPath, secret.Data)
	return err
}

func (v *vaultSecretStorageClient) secretPath(project string) (string, error) {
	var tpl bytes.Buffer
	data := struct {
		Project string
	}{
		Project: project,
	}

	err := v.secretPathTemplate.Execute(&tpl, data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (v *vaultSecretStorageClient) login() error {
	// TODO: properly authenticate to Vault
	v.vaultClient.SetToken(v.vaultConfig.Token)
	return nil
}
