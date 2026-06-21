package services

import (
	"database/sql"
	"errors"

	"DNSPilot/internal/models"
)

type SettingsService struct {
	db *sql.DB
}

func NewSettingsService(db *sql.DB) *SettingsService {
	return &SettingsService{db: db}
}

func defaultSettings() models.Settings {
	return models.Settings{
		Locale:             "fa",
		Theme:              "dark",
		AutoApplyFastest:   false,
		BenchmarkAttempts:  5,
		BenchmarkTimeoutMS: 1500,
		TestDomain:         "google.com",
		SelectedAdapterID:  "",
		LastProfileID:      "",
	}
}

func (s *SettingsService) GetSettings() (models.Settings, error) {
	var out models.Settings

	row := s.db.QueryRow(`
		SELECT locale, theme, auto_apply_fastest, benchmark_attempts, benchmark_timeout_ms, test_domain, selected_adapter_id, last_profile_id
		FROM settings
		WHERE id = 1
	`)

	err := row.Scan(
		&out.Locale,
		&out.Theme,
		&out.AutoApplyFastest,
		&out.BenchmarkAttempts,
		&out.BenchmarkTimeoutMS,
		&out.TestDomain,
		&out.SelectedAdapterID,
		&out.LastProfileID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return defaultSettings(), nil
	}
	if err != nil {
		return models.Settings{}, err
	}

	return out, nil
}

func (s *SettingsService) SaveSettings(settings models.Settings) error {
	if settings.Locale == "" {
		settings.Locale = "fa"
	}
	if settings.Theme == "" {
		settings.Theme = "dark"
	}
	if settings.BenchmarkAttempts <= 0 {
		settings.BenchmarkAttempts = 5
	}
	if settings.BenchmarkTimeoutMS <= 0 {
		settings.BenchmarkTimeoutMS = 1500
	}
	if settings.TestDomain == "" {
		settings.TestDomain = "google.com"
	}

	_, err := s.db.Exec(`
		INSERT INTO settings (
			id, locale, theme, auto_apply_fastest, benchmark_attempts, benchmark_timeout_ms, test_domain, selected_adapter_id, last_profile_id
		) VALUES (
			1, ?, ?, ?, ?, ?, ?, ?, ?
		)
		ON CONFLICT(id) DO UPDATE SET
			locale = excluded.locale,
			theme = excluded.theme,
			auto_apply_fastest = excluded.auto_apply_fastest,
			benchmark_attempts = excluded.benchmark_attempts,
			benchmark_timeout_ms = excluded.benchmark_timeout_ms,
			test_domain = excluded.test_domain,
			selected_adapter_id = excluded.selected_adapter_id,
			last_profile_id = excluded.last_profile_id
	`,
		settings.Locale,
		settings.Theme,
		settings.AutoApplyFastest,
		settings.BenchmarkAttempts,
		settings.BenchmarkTimeoutMS,
		settings.TestDomain,
		settings.SelectedAdapterID,
		settings.LastProfileID,
	)
	return err
}
