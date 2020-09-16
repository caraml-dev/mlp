// +build integration

package database

import "os"

var (
	host     = "postgres"
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	database = os.Getenv("POSTGRES_DB")
)
