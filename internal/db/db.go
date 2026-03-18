package db

import "database/sql"

// DB is the common interface every database driver must implement
type DB interface {
	Connect() error
	Close() error
	Tables() ([]string, error)
	Columns(table string) ([]Column, error)
	Rows(table string, limit int) (*sql.Rows, error)
	Query(query string) (*sql.Rows, error)
	CountRows(table string) (int, error)
}

// Column holds metadata about a single table column
type Column struct {
	Name     string
	Type     string
	Nullable bool
}

// Config holds the connection details
type Config struct {
	Driver   string // postgres, mysql, sqlite
	DSN      string // full connection string
}

// New returns the correct driver based on config
func New(cfg Config) DB {
	switch cfg.Driver {
	case "postgres":
		return &Postgres{dsn: cfg.DSN}
	case "mysql":
		return &MySQL{dsn: cfg.DSN}
	case "sqlite":
		return &SQLite{dsn: cfg.DSN}
	default:
		return nil
	}
}