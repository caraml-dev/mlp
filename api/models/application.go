package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Application struct {
	Id          Id                 `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Href        string             `json:"href"`
	IconName    string             `json:"icon" gorm:"column:icon"`
	UseProjects bool               `json:"use_projects"`
	IsInBeta    bool               `json:"is_in_beta"`
	IsDisabled  bool               `json:"is_disabled"`
	Config      *ApplicationConfig `json:"config"`
}

type ApplicationConfig struct {
	Sections []ApplicationSection `json:"sections"`
}

func (c ApplicationConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ApplicationConfig) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type ApplicationSection struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

func (c ApplicationSection) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ApplicationSection) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}
