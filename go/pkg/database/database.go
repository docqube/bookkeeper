package database

import (
	"database/sql"
	"flag"
	"fmt"
	"sync"
	"time"
)

var (
	username         = flag.String("db-user", "postgres", "Database username")
	password         = flag.String("db-password", "postgres", "Database password")
	host             = flag.String("db-host", "localhost", "Database Host")
	databaseName     = flag.String("db-name", "bookkeeper", "Database Name")
	sqlSSL           = flag.Bool("db-ssl", false, "Sets sslmode=disable if false")
	databaseInstance = make(map[string]*sql.DB, 0)
	lock             sync.Mutex
	// PGMaxConnection is the amount of connections that the postgres pool can hold
	PGMaxConnection = 5
)

// GetDBWithName returns a *sql.DB to the given database
func GetConnection() (*sql.DB, error) {
	lock.Lock()
	defer lock.Unlock()
	if database, ok := databaseInstance[*databaseName]; ok {
		return database, nil
	}
	database, err := sql.Open("postgres", DBConnectionString(*databaseName))
	// If there was no error, configure and cache the database connection
	if err == nil {
		database.SetMaxOpenConns(PGMaxConnection)
		database.SetMaxIdleConns(PGMaxConnection / 2)
		databaseInstance[*databaseName] = database
	}
	return database, err
}

func DBConnectionString(databaseName string) string {
	datasource := fmt.Sprintf("postgres://%s:%s@%s/%s?timezone=UTC", *username, *password, *host, databaseName)
	if !*sqlSSL {
		datasource += "&sslmode=disable"
	}
	return datasource
}

// NormalizeTime truncates the Milliseconds and convert the time to UTC
func NormalizeTime(timestamp time.Time) time.Time {
	return timestamp.Truncate(time.Millisecond).UTC()
}
