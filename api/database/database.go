package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jinzhu/gorm"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	// required for gomigrate
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// enable postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/caraml-dev/mlp/api/config"
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
		return nil, err
	}

	// Get the underlying SQL DB and apply connection properties
	sqlDB := db.DB()
	if sqlDB == nil {
		return nil, errors.New("Failed to get underlying database connection")
	}
	sqlDB.SetConnMaxIdleTime(dbCfg.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(dbCfg.ConnMaxLifetime)
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)

	db.LogMode(false)
	err = runDBMigration(db, dbCfg.MigrationPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runDBMigration(db *gorm.DB, migrationPath string) error {
	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
