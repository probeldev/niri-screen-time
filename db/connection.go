package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

type DBConnection struct {
	db    *sql.DB
	sem   chan struct{} // Семафор для ограничения параллелизма
	mutex sync.Mutex    // Мьютекс для операций DDL (CREATE TABLE и т.д.)
}

func NewDBConnection() (*DBConnection, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL&_mutex=full", dbPath)

	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ограничиваем количество одновременных операций записи
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL;"); err != nil {
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	return &DBConnection{
		db:  db,
		sem: make(chan struct{}, 1), // Буфер = 1 (только одна операция записи за раз)
	}, nil
}

// Close закрывает подключение к БД
func (dbc *DBConnection) Close() error {
	close(dbc.sem)
	return dbc.db.Close()
}

// InitTables создает таблицы (использует мьютекс для безопасности)
func (dbc *DBConnection) InitTables() error {
	dbc.mutex.Lock()
	defer dbc.mutex.Unlock()

	_, err := dbc.db.Exec(`
	CREATE TABLE IF NOT EXISTS screen_time (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TIMESTAMP NOT NULL,
		app_id TEXT NOT NULL,
		title TEXT NOT NULL,
		sleep INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS aggregated_screen_time (
		date TIMESTAMP NOT NULL,
		app_id TEXT NOT NULL,
		title TEXT NOT NULL,
		sleep INTEGER NOT NULL
	);
	`)
	return err
}

// Exec выполняет запрос с ограничением параллелизма
func (dbc *DBConnection) Exec(query string, args ...any) (sql.Result, error) {
	dbc.sem <- struct{}{}        // Захватываем слот
	defer func() { <-dbc.sem }() // Освобождаем

	return dbc.db.Exec(query, args...)
}

// Query выполняет запрос на чтение (без ограничений)
func (dbc *DBConnection) Query(query string, args ...any) (*sql.Rows, error) {
	return dbc.db.Query(query, args...)
}

// GetDB возвращает *sql.DB (для обратной совместимости)
func (dbc *DBConnection) GetDB() *sql.DB {
	return dbc.db
}

func (dbc *DBConnection) Vacuum() error {
	_, err := dbc.db.Exec("VACUUM")
	return err
}
