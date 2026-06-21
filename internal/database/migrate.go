package database

import "database/sql"

func Migrate(db *sql.DB) error {

	query := `
	CREATE TABLE IF NOT EXISTS custom_dns (
		id TEXT PRIMARY KEY,
		name TEXT,
		primary_ip TEXT,
		secondary_ip TEXT
	);
	`

	_, err := db.Exec(query)

	return err
}
