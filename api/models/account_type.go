package models

type AccountType string

var AccountTypes = struct {
	Gitlab AccountType
}{
	Gitlab: "Gitlab",
}
