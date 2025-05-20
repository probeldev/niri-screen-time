package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/probeldev/niri-screen-time/model"

	_ "github.com/mattn/go-sqlite3"
)

type ScreenTimeDB struct {
	db *sql.DB
}

// NewScreenTimeDB инициализирует подключение к БД и создает таблицу если нужно
func NewScreenTimeDB(dbPath string) (*ScreenTimeDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Включаем WAL режим для лучшей производительности при конкурентном доступе
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Создаем таблицу с оптимальными индексами
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS screen_time (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TIMESTAMP NOT NULL,
		app_id TEXT NOT NULL,
		title TEXT NOT NULL,
		sleep INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_date ON screen_time(date);
	CREATE INDEX IF NOT EXISTS idx_app_id ON screen_time(app_id);
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
