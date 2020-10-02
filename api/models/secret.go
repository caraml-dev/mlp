package models

import "github.com/gojek/mlp/api/util"

type Secret struct {
	Id        Id     `json:"id"`
	ProjectId Id     `json:"project_id"`
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

func (s *Secret) isValid(checkingId bool) bool {
	if checkingId && s.Id <= 0 {
		return false
	}
	maxNameChar := 100
	if s.Name == "" || len(s.Name) > maxNameChar {
		return false
	}
	if s.Data == "" {
		return false
	}
	if s.ProjectId == 0 {
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
		Id:             s.Id,
		ProjectId:      s.ProjectId,
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
		Id:             s.Id,
		ProjectId:      s.ProjectId,
		Name:           s.Name,
		Data:           encryptedData,
		CreatedUpdated: s.CreatedUpdated,
	}, nil
}
