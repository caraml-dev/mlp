package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jinzhu/gorm"

	// required for gomigrate
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// enable postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/gojek/mlp/api/config"
)

// InitDB initialises a database connection as well as runs the migration scripts.
// It is important to close the database after using it by calling defer db.Close()
func InitDB(dbCfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.User,
			dbCfg.Database,
			dbCfg.Password))
	if err != nil {
		panic(err)
	}
	db.LogMode(false)
	runDBMigration(db, dbCfg.MigrationPath)

	return db, nil
}

func runDBMigration(db *gorm.DB, migrationPath string) {
	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
