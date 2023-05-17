package models

import "github.com/caraml-dev/mlp/api/util"

// Secret represents user defined secret
type Secret struct {
	// ID is the unique identifier of the secret
	ID ID `json:"id"`
	// ProjectID is the unique identifier of the project
	ProjectID ID `json:"project_id"`
	// Project is the project of the secret
	Project *Project `json:"-"`
	// Name is the name of the secret
	Name string `json:"name"`
	// Data is secret value
	Data string `json:"data"`
	// SecretStorageID is the unique identifier of the secret storage for storing the secret
	SecretStorageID *ID `json:"secret_storage_id,omitempty"`
	// SecretStorage is the secret storage for storing the secret
	SecretStorage *SecretStorage `json:"secret_storage,omitempty"`
	// CreatedUpdated is the timestamp of the secret creation and update
	CreatedUpdated
}

func (s *Secret) IsValidForInsertion() bool {
	return s.isValid(false)
}

func (s *Secret) IsValidForMutation() bool {
	return s.isValid(true)
}

func (s *Secret) isValid(checkingID bool) bool {
	if checkingID && s.ID <= 0 {
		return false
	}
	maxNameChar := 100
	if s.Name == "" || len(s.Name) > maxNameChar {
		return false
	}
	if s.Data == "" {
		return false
	}
	if s.ProjectID == 0 {
		return false
	}
	return true
}

func (s *Secret) CopyValueFrom(secret *Secret) {
	if secret.Name != "" {
		s.Name = secret.Name
	}
	if secret.Data != "" {
		s.Data = secret.Data
	}
}

func (s *Secret) DecryptData(passphrase string) (*Secret, error) {
	encryptedData := s.Data
	decryptedData, err := util.Decrypt(encryptedData, passphrase)
	if err != nil {
		return nil, err
	}

	return &Secret{
		ID:             s.ID,
		ProjectID:      s.ProjectID,
		Name:           s.Name,
		Data:           decryptedData,
		CreatedUpdated: s.CreatedUpdated,
	}, nil
}

func (s *Secret) EncryptData(passphrase string) (*Secret, error) {
	plainText := s.Data
	encryptedData, err := util.Encrypt(plainText, passphrase)
	if err != nil {
		return nil, err
	}
	return &Secret{
		ID:             s.ID,
		ProjectID:      s.ProjectID,
		Name:           s.Name,
		Data:           encryptedData,
		CreatedUpdated: s.CreatedUpdated,
	}, nil
}
