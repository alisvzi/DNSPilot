package models

type Settings struct {
	Locale string `json:"locale"`
	Theme  string `json:"theme"`

	AutoApplyFastest bool `json:"auto_apply_fastest"`

	BenchmarkAttempts  int    `json:"benchmark_attempts"`
	BenchmarkTimeoutMS int    `json:"benchmark_timeout_ms"`
	TestDomain         string `json:"test_domain"`

	SelectedAdapterID string `json:"selected_adapter_id"`
	LastProfileID     string `json:"last_profile_id"`
}
