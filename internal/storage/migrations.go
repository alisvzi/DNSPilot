package storage

import "database/sql"

func Migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			locale TEXT NOT NULL DEFAULT 'fa',
			theme TEXT NOT NULL DEFAULT 'dark',
			auto_apply_fastest INTEGER NOT NULL DEFAULT 0,
			benchmark_attempts INTEGER NOT NULL DEFAULT 5,
			benchmark_timeout_ms INTEGER NOT NULL DEFAULT 1500,
			test_domain TEXT NOT NULL DEFAULT 'google.com',
			selected_adapter_id TEXT NOT NULL DEFAULT '',
			last_profile_id TEXT NOT NULL DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS custom_dns_lists (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS custom_dns_servers (
			id TEXT PRIMARY KEY,
			list_id TEXT NOT NULL,
			name TEXT NOT NULL,
			primary_ipv4 TEXT NOT NULL DEFAULT '',
			secondary_ipv4 TEXT NOT NULL DEFAULT '',
			primary_ipv6 TEXT NOT NULL DEFAULT '',
			secondary_ipv6 TEXT NOT NULL DEFAULT '',
			provider TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			category TEXT NOT NULL DEFAULT 'custom',
			tags_json TEXT NOT NULL DEFAULT '[]',
			is_custom INTEGER NOT NULL DEFAULT 1,
			enabled INTEGER NOT NULL DEFAULT 1,
			FOREIGN KEY(list_id) REFERENCES custom_dns_lists(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS dns_backups (
			adapter_name TEXT PRIMARY KEY,
			servers_json TEXT NOT NULL,
			is_ipv6 INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS benchmark_profiles (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			attempts INTEGER NOT NULL DEFAULT 5,
			timeout_ms INTEGER NOT NULL DEFAULT 1500
		);`,
		`CREATE TABLE IF NOT EXISTS benchmark_profile_domains (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id TEXT NOT NULL,
			domain TEXT NOT NULL,
			FOREIGN KEY(profile_id) REFERENCES benchmark_profiles(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS benchmark_results (
			id TEXT PRIMARY KEY,
			server_id TEXT NOT NULL,
			server_name TEXT NOT NULL,
			profile_id TEXT NOT NULL,
			profile_name TEXT NOT NULL,
			latency_ms REAL NOT NULL,
			jitter_ms REAL NOT NULL,
			success_rate REAL NOT NULL,
			packet_loss REAL NOT NULL,
			score REAL NOT NULL,
			attempts INTEGER NOT NULL,
			created_at TEXT NOT NULL
		);`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	_, err := db.Exec(`
		INSERT OR IGNORE INTO settings (
			id, locale, theme, auto_apply_fastest, benchmark_attempts, benchmark_timeout_ms, test_domain, selected_adapter_id, last_profile_id
		) VALUES (
			1, 'fa', 'dark', 0, 5, 1500, 'google.com', '', ''
		);
	`)
	return err
}
