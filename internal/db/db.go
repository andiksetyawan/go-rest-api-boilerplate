package db

import "database/sql"

type DB interface {
	DSN() string
	Connect() (*sql.DB, error)
	MigrateUp() error
	MigrateDown() error
}
