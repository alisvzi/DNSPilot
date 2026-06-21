package models

type BenchmarkResult struct {
	ID string `json:"id"`

	ServerID   string `json:"server_id"`
	ServerName string `json:"server_name"`

	ProfileID   string `json:"profile_id"`
	ProfileName string `json:"profile_name"`

	LatencyMs   float64 `json:"latency_ms"`
	JitterMs    float64 `json:"jitter_ms"`
	SuccessRate float64 `json:"success_rate"`
	PacketLoss  float64 `json:"packet_loss"`
	Score       float64 `json:"score"`
	Attempts    int     `json:"attempts"`

	CreatedAt string `json:"created_at"`
}
