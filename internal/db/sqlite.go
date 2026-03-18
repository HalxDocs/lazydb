package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	dsn  string
	conn *sql.DB
}

func (s *SQLite) Connect() error {
	conn, err := sql.Open("sqlite3", s.dsn)
	if err != nil {
		return fmt.Errorf("sqlite connect: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("sqlite ping: %w", err)
	}
	s.conn = conn
	return nil
}

func (s *SQLite) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *SQLite) Tables() ([]string, error) {
	rows, err := s.conn.Query(`
		SELECT name FROM sqlite_master
		WHERE type='table'
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("sqlite tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, nil
}

func (s *SQLite) Columns(table string) ([]Column, error) {
	rows, err := s.conn.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return nil, fmt.Errorf("sqlite columns: %w", err)
	}
	defer rows.Close()

	var cols []Column
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull int
		var dfltValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return nil, err
		}
		cols = append(cols, Column{
			Name:     name,
			Type:     colType,
			Nullable: notNull == 0,
		})
	}
	return cols, nil
}

func (s *SQLite) Rows(table string, limit int) (*sql.Rows, error) {
	return s.conn.Query(fmt.Sprintf("SELECT * FROM %s LIMIT ?", table), limit)
}

func (s *SQLite) Query(query string) (*sql.Rows, error) {
	return s.conn.Query(query)
}