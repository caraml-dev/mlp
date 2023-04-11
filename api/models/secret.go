package models

import "github.com/caraml-dev/mlp/api/util"

type Secret struct {
	ID        ID     `json:"id"`
	ProjectID ID     `json:"project_id"`
	Name      string `json:"name"`
	Data      string `json:"data"`
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
