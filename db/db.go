package db

import (
	"fmt"
	"os"
	"path/filepath"
)

// getDbPath возвращает путь к файлу базы данных
func getDbPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".local", "share", "niri-screen-time")
	dbPath := filepath.Join(dbDir, "db.db")

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dbDir, err)
	}

	return dbPath, nil
}
