package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// DBMigration performs a database migration of the database using the supplied files.
// It returns nil if nothing changed or migration was successful and an error otherwise.
func DBMigration(migrationFilesPath string) error {
	migration, err := migrate.New(migrationFilesPath, DBConnectionString(*databaseName))
	if err != nil {
		return fmt.Errorf("connecting to database: %s", err)
	}

	err = migration.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return fmt.Errorf("running migration: %s", err)
	}
	return nil
}
