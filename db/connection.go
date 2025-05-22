package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DBConnection struct {
	db *sql.DB
}

func NewDBConnection() (*DBConnection, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL;"); err != nil {
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	return &DBConnection{db: db}, nil
}

func (dbc *DBConnection) Close() error {
	return dbc.db.Close()
}

func (dbc *DBConnection) InitTables() error {
	_, err := dbc.db.Exec(`
		CREATE TABLE IF NOT EXISTS screen_time (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TIMESTAMP NOT NULL,
			app_id TEXT NOT NULL,
			title TEXT NOT NULL,
			sleep INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS aggregated_screen_time (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TIMESTAMP NOT NULL,
			app_id TEXT NOT NULL,
			title TEXT NOT NULL,
			sleep INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func (dbc *DBConnection) GetDB() *sql.DB {
	return dbc.db
}
