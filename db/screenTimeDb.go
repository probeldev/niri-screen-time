package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/probeldev/niri-screen-time/model"

	_ "github.com/mattn/go-sqlite3"
)

type ScreenTimeDB struct {
	db *sql.DB
}

func getDbPath() (string, error) {
	// Раскрываем ~ в домашнюю директорию
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Формируем полный путь
	dbDir := filepath.Join(homeDir, ".local", "share", "niri-screen-time")
	dbPath := filepath.Join(dbDir, "db.db")

	// Создаём директорию (если её нет)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dbDir, err)
	}

	return dbPath, nil
}

// NewScreenTimeDB инициализирует подключение к БД и создаёт таблицы.
func NewScreenTimeDB() (*ScreenTimeDB, error) {
	// Получаем путь к БД
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	// Открываем/создаём БД
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройки SQLite
	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL;"); err != nil {
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS screen_time (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TIMESTAMP NOT NULL,
		app_id TEXT NOT NULL,
		title TEXT NOT NULL,
		sleep INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &ScreenTimeDB{db: db}, nil
}

// Close закрывает соединение с БД
func (stdb *ScreenTimeDB) Close() error {
	return stdb.db.Close()
}

// Insert добавляет новую запись
func (stdb *ScreenTimeDB) Insert(st model.ScreenTime) error {
	_, err := stdb.db.Exec(
		"INSERT INTO screen_time(date, app_id, title, sleep) VALUES(?, ?, ?, ?)",
		st.Date, st.AppID, st.Title, st.Sleep,
	)
	return err
}

// BulkInsert добавляет несколько записей в одной транзакции
func (stdb *ScreenTimeDB) BulkInsert(records []model.ScreenTime) error {
	tx, err := stdb.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO screen_time(date, app_id, title, sleep) VALUES(?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, st := range records {
		if _, err := stmt.Exec(st.Date, st.AppID, st.Title, st.Sleep); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetByDateRange возвращает записи за указанный период
func (stdb *ScreenTimeDB) GetByDateRange(from, to time.Time) ([]model.ScreenTime, error) {
	rows, err := stdb.db.Query(
		"SELECT date, app_id, title, sleep FROM screen_time WHERE date BETWEEN ? AND ? ORDER BY date",
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ScreenTime
	for rows.Next() {
		var st model.ScreenTime
		if err := rows.Scan(&st.Date, &st.AppID, &st.Title, &st.Sleep); err != nil {
			return nil, err
		}
		results = append(results, st)
	}

	return results, nil
}

// GetAppUsage возвращает суммарное время по приложениям за период
func (stdb *ScreenTimeDB) GetAppUsage(from, to time.Time) (map[string]int, error) {
	rows, err := stdb.db.Query(
		"SELECT app_id, SUM(sleep) FROM screen_time WHERE date BETWEEN ? AND ? GROUP BY app_id",
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var appID string
		var total int
		if err := rows.Scan(&appID, &total); err != nil {
			return nil, err
		}
		result[appID] = total
	}

	return result, nil
}

// DeleteOldRecords удаляет записи старше указанной даты
func (stdb *ScreenTimeDB) DeleteOldRecords(before time.Time) error {
	_, err := stdb.db.Exec("DELETE FROM screen_time WHERE date < ?", before)
	return err
}
