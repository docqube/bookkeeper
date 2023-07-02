package database

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"docqube.de/bookkeeper/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrateFS embed.FS

const PGMaxConnections = 5

// buildConnectionString returns a valid PostgreSQL connection string
// with the passed database configuration.
func buildConnectionString(config config.DatabaseConfig) string {
	datasource := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?timezone=UTC",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
	if config.SSLMode != "" {
		datasource += "&sslmode=disable"
	}
	return datasource
}

// InitializeDatabase initializes the database connection and performs a database migration.
// After a successful initialization, the database connection is returned.
func InitializeDatabase(config *config.Config) (*sql.DB, error) {
	database, err := sql.Open("postgres", buildConnectionString(config.DatabaseConfig))
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(PGMaxConnections)
	database.SetMaxIdleConns(PGMaxConnections / 2)

	err = migrateDatabase(database)
	if err != nil {
		return nil, err
	}
	return database, nil
}

// NormalizeTime truncates the Milliseconds and convert the time to UTC
func NormalizeTime(timestamp time.Time) time.Time {
	return timestamp.Truncate(time.Millisecond).UTC()
}

// migrateDatabase performs a database migration with the passed database connection.
// The migration files are embedded and will be executed in the order of their filenames.
func migrateDatabase(db *sql.DB) error {
	sourceInstance, err := iofs.New(migrateFS, "migrations")
	if err != nil {
		return err
	}

	databaseDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithInstance(
		"iofs",
		sourceInstance,
		"postgres",
		databaseDriver,
	)
	if err != nil {
		return err
	}

	err = migration.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}
	return nil
}
