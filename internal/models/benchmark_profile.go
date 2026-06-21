package models

type BenchmarkProfile struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Description string   `json:"description"`
	Domains     []string `json:"domains"`

	Attempts  int `json:"attempts"`
	TimeoutMS int `json:"timeout_ms"`
}
