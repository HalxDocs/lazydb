package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	dsn  string
	conn *sql.DB
}

func (p *Postgres) Connect() error {
	conn, err := sql.Open("postgres", p.dsn)
	if err != nil {
		return fmt.Errorf("postgres connect: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("postgres ping: %w", err)
	}
	p.conn = conn
	return nil
}

func (p *Postgres) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *Postgres) Tables() ([]string, error) {
	rows, err := p.conn.Query(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`)
	if err != nil {
		return nil, fmt.Errorf("postgres tables: %w", err)
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

func (p *Postgres) Columns(table string) ([]Column, error) {
	rows, err := p.conn.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`, table)
	if err != nil {
		return nil, fmt.Errorf("postgres columns: %w", err)
	}
	defer rows.Close()

	var cols []Column
	for rows.Next() {
		var c Column
		var nullable string
		if err := rows.Scan(&c.Name, &c.Type, &nullable); err != nil {
			return nil, err
		}
		c.Nullable = nullable == "YES"
		cols = append(cols, c)
	}
	return cols, nil
}

func (p *Postgres) Rows(table string, limit int) (*sql.Rows, error) {
	return p.conn.Query(fmt.Sprintf("SELECT * FROM %s LIMIT $1", table), limit)
}

func (p *Postgres) Query(query string) (*sql.Rows, error) {
	return p.conn.Query(query)
}

func (p *Postgres) CountRows(table string) (int, error) {
	var count int
	err := p.conn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
	return count, err
}