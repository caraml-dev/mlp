package models

import (
	"strconv"
	"time"
)

type ID int

func (id ID) String() string {
	return strconv.Itoa(int(id))
}

func ParseID(id string) (ID, error) {
	parsed, err := strconv.Atoi(id)
	if err != nil {
		return -1, err
	}
	return ID(parsed), nil
}

type CreatedUpdated struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
