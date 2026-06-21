package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(base, "DNSPilot")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dir, "dnspilot.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
