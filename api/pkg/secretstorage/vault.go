package secretstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	vault "github.com/hashicorp/vault/api"

	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/models"
	mlperror "github.com/caraml-dev/mlp/api/pkg/errors"
)

type vaultSecretStorageClient struct {
	secretPathTemplate *template.Template
	vaultClient        *vault.Client
	vaultConfig        *models.VaultConfig
	authHelper         authHelper
}

// NewVaultSecretStorageClient creates a new secret storage client backed by Vault
func NewVaultSecretStorageClient(ss *models.SecretStorage) (Client, error) {
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

	// Authenticate to Vault
	switch ss.Config.VaultConfig.AuthMethod {
	case models.TokenAuthMethod:
		vaultClient.SetToken(ss.Config.VaultConfig.Token)
	case models.GCPAuthMethod:
		cli.authHelper = newGcpAuthHelper(ss.Config.VaultConfig)
		_, err := cli.authHelper.login(vaultClient)
		if err != nil {
			return nil, fmt.Errorf("failed to login to vault: %w", err)
		}
		// run token renewal in background
		go cli.renewToken()
	default:
		return nil, fmt.Errorf("unknown auth method: %s", ss.Config.VaultConfig.AuthMethod)
	}

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

// Set creates or updates a CaraML secret of a project in Vault
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

// List lists all CaraML secrets of a project in Vault
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

// Delete deletes a CaraML secret of a project in Vault
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

	delete(secret.Data, name)

	_, err = v.vaultClient.KVv2(v.vaultConfig.MountPath).Put(context.Background(), secretPath, secret.Data)
	return err
}

func (v *vaultSecretStorageClient) DeleteAll(project string) error {
	secretPath, err := v.secretPath(project)
	if err != nil {
		return err
	}

	return v.vaultClient.KVv2(v.vaultConfig.MountPath).Delete(context.Background(), secretPath)
}

func (v *vaultSecretStorageClient) SetAll(secrets map[string]string, project string) error {
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

	for k, v := range secrets {
		existingSecrets[k] = v
	}

	_, err = v.vaultClient.KVv2(v.vaultConfig.MountPath).Put(context.Background(), secretPath, existingSecrets)
	return err
}

// secretPath returns the secret path for a project
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

// renewToken renews the Vault token periodically or attempt to do login if the token is not renewable
// adapted from https://github.com/hashicorp/vault-examples/blob/main/examples/token-renewal/go/example.go
func (v *vaultSecretStorageClient) renewToken() {
	for {
		vaultLoginResp, err := v.authHelper.login(v.vaultClient)
		if err != nil {
			log.Errorf("unable to authenticate to Vault: %v", err)
			continue
		}
		tokenErr := v.manageTokenLifecycle(vaultLoginResp)
		if tokenErr != nil {
			log.Errorf("unable to start managing token lifecycle: %v", tokenErr)
			continue
		}
	}
}

// Starts token lifecycle management. Returns only fatal errors as errors,
// otherwise returns nil, so we can attempt login again.
func (v *vaultSecretStorageClient) manageTokenLifecycle(token *vault.Secret) error {
	renew := token.Auth.Renewable
	if !renew {
		log.Infof("Token is not configured to be renewable. Re-attempting login.")
		return nil
	}

	watcher, err := v.vaultClient.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
		Secret: token,
	})
	if err != nil {
		return fmt.Errorf("unable to initialize new lifetime watcher for renewing auth token: %w", err)
	}

	go watcher.Start()
	defer watcher.Stop()

	for {
		select {
		// `DoneCh` will return if renewal fails, or if the remaining lease
		// duration is under a built-in threshold and either renewing is not
		// extending it or renewing is disabled. In any case, the caller
		// needs to attempt to log in again.
		case err := <-watcher.DoneCh():
			if err != nil {
				log.Infof("Failed to renew token: %v. Re-attempting login.", err)
				return nil
			}
			// This occurs once the token has reached max TTL.
			log.Infof("Token can no longer be renewed. Re-attempting login.")
			return nil

		// Successfully completed renewal
		case renewal := <-watcher.RenewCh():
			log.Infof("Successfully renewed: %#v", renewal)
		}
	}
}
