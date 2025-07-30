// Package db implements SQLite storage.
package db

import (
	"fmt"
	"os"
	"path/filepath"
)

// getDBPath возвращает путь к файлу базы данных
func getDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".local", "share", "niri-screen-time")
	dbPath := filepath.Join(dbDir, "db.db")

	var perm uint32 = 0755

	if err := os.MkdirAll(dbDir, os.FileMode(perm)); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dbDir, err)
	}

	return dbPath, nil
}
