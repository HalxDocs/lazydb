package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	dsn  string
	conn *sql.DB
}

func (m *MySQL) Connect() error {
	conn, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return fmt.Errorf("mysql connect: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("mysql ping: %w", err)
	}
	m.conn = conn
	return nil
}

func (m *MySQL) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

func (m *MySQL) Tables() ([]string, error) {
	rows, err := m.conn.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("mysql tables: %w", err)
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

func (m *MySQL) Columns(table string) ([]Column, error) {
	rows, err := m.conn.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_name = ?
		ORDER BY ordinal_position
	`, table)
	if err != nil {
		return nil, fmt.Errorf("mysql columns: %w", err)
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

func (m *MySQL) Rows(table string, limit int) (*sql.Rows, error) {
	return m.conn.Query(fmt.Sprintf("SELECT * FROM %s LIMIT ?", table), limit)
}

func (m *MySQL) Query(query string) (*sql.Rows, error) {
	return m.conn.Query(query)
}

func (m *MySQL) CountRows(table string) (int, error) {
	var count int
	err := m.conn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
	return count, err
}