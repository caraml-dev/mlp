//go:build integration

package database

import "os"

var (
	host     = os.Getenv("DATABASE_HOST")
	user     = os.Getenv("DATABASE_USER")
	password = os.Getenv("DATABASE_PASSWORD")
	database = os.Getenv("DATABASE_NAME")
)
