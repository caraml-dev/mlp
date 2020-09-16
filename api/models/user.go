package models

import (
	"strconv"
	"time"
)

type Id int

func (id Id) String() string {
	return strconv.Itoa(int(id))
}

func ParseId(id string) (Id, error) {
	if parsed, err := strconv.Atoi(id); err != nil {
		return -1, err
	} else {
		return Id(parsed), nil
	}
}

type CreatedUpdated struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	Id       Id     `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
	CreatedUpdated
}
