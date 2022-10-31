package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/lib/pq"
)

type Project struct {
	ID   ID     `json:"id"`
	Name string `json:"name" validate:"required,min=3,max=50,subdomain_rfc1123"`
	// nolint:lll // Next line is 121 characters (lll)
	MLFlowTrackingURL string         `json:"mlflow_tracking_url" gorm:"column:mlflow_tracking_url" validate:"omitempty,url"`
	Administrators    pq.StringArray `json:"administrators" gorm:"column:administrators;type:varchar(256)[]"`
	Readers           pq.StringArray `json:"readers" gorm:"column:readers;type:varchar(256)[]"`
	Team              string         `json:"team" validate:"required,min=1,max=64"`
	Stream            string         `json:"stream" validate:"required,min=1,max=64"`
	Labels            Labels         `json:"labels,omitempty" gorm:"column:labels"`
	CreatedUpdated
}

type Labels []Label

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (labels Labels) Value() (driver.Value, error) {
	return json.Marshal(labels)
}

func (labels *Labels) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &labels)
}
