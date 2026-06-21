package benchmark

type Result struct {
	ServerIP  string `json:"server_ip"`
	LatencyMS int64  `json:"latency_ms"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}
