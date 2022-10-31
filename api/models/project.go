package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/lib/pq"
)

type Project struct {
	ID                ID             `json:"id"`
	Name              string         `json:"name" validate:"required,min=3,max=50,subdomain_rfc1123"`
	MLFlowTrackingURL string         `json:"mlflow_tracking_url" gorm:"mlflow_tracking_url" validate:"omitempty,url"`
	Administrators    pq.StringArray `json:"administrators" gorm:"administrators;type:varchar(256)[]"`
	Readers           pq.StringArray `json:"readers" gorm:"readers;type:varchar(256)[]"`
	Team              string         `json:"team" validate:"required,min=1,max=64"`
	Stream            string         `json:"stream" validate:"required,min=1,max=64"`
	Labels            Labels         `json:"labels,omitempty" gorm:"labels"`
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
