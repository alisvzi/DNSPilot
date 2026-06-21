package benchmark

import (
	"fmt"
	"math"
	"net"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"DNSPilot/internal/models"

	"github.com/miekg/dns"
)

func Run(profile models.BenchmarkProfile, servers []string) ([]models.BenchmarkResult, error) {
	if len(servers) == 0 {
		return nil, fmt.Errorf("no DNS servers provided")
	}

	if profile.Attempts <= 0 {
		profile.Attempts = 5
	}
	if profile.TimeoutMS <= 0 {
		profile.TimeoutMS = 1500
	}
	if len(profile.Domains) == 0 {
		profile.Domains = []string{"google.com"}
	}

	workers := runtime.NumCPU() * 4
	if workers < 1 {
		workers = 1
	}
	if workers > len(servers) {
		workers = len(servers)
	}

	jobs := make(chan string)
	results := make(chan models.BenchmarkResult, len(servers))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for server := range jobs {
				results <- benchmarkServer(server, profile)
			}
		}()
	}

	go func() {
		for _, server := range servers {
			server = strings.TrimSpace(server)
			if server != "" {
				jobs <- server
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]models.BenchmarkResult, 0, len(servers))
	for r := range results {
		out = append(out, r)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].LatencyMs < out[j].LatencyMs
		}
		return out[i].Score < out[j].Score
	})

	return out, nil
}

func benchmarkServer(server string, profile models.BenchmarkProfile) models.BenchmarkResult {
	attempts := profile.Attempts
	if attempts <= 0 {
		attempts = 5
	}

	timeoutMS := profile.TimeoutMS
	if timeoutMS <= 0 {
		timeoutMS = 1500
	}

	totalQueries := attempts * len(profile.Domains)
	if totalQueries == 0 {
		totalQueries = attempts
	}

	var durations []float64
	successes := 0

	for i := 0; i < attempts; i++ {
		for _, domain := range profile.Domains {
			d, err := queryOnce(server, domain, time.Duration(timeoutMS)*time.Millisecond)
			if err == nil {
				durations = append(durations, float64(d.Milliseconds()))
				successes++
			}
		}
	}

	latency := float64(timeoutMS)
	jitter := 0.0

	if len(durations) > 0 {
		latency = medianFloat(durations)
		jitter = stddevFloat(durations)
	}

	successRate := (float64(successes) / float64(totalQueries)) * 100
	packetLoss := 100 - successRate
	score := (latency * 0.7) + (jitter * 0.2) + (packetLoss * 2)

	return models.BenchmarkResult{
		ID:          newBenchmarkID(server),
		ServerID:    server,
		ServerName:  server,
		ProfileID:   profile.ID,
		ProfileName: profile.Name,
		LatencyMs:   round2(latency),
		JitterMs:    round2(jitter),
		SuccessRate: round2(successRate),
		PacketLoss:  round2(packetLoss),
		Score:       round2(score),
		Attempts:    totalQueries,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func queryOnce(server, domain string, timeout time.Duration) (time.Duration, error) {
	client := &dns.Client{
		Net:     "udp",
		Timeout: timeout,
	}

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	start := time.Now()
	_, _, err := client.Exchange(msg, net.JoinHostPort(server, "53"))
	if err == nil {
		return time.Since(start), nil
	}

	client.Net = "tcp"
	start = time.Now()
	_, _, err = client.Exchange(msg, net.JoinHostPort(server, "53"))
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func medianFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	cp := append([]float64(nil), values...)
	sort.Float64s(cp)

	mid := len(cp) / 2
	if len(cp)%2 == 1 {
		return cp[mid]
	}
	return (cp[mid-1] + cp[mid]) / 2
}

func stddevFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	var variance float64
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))

	return math.Sqrt(variance)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func newBenchmarkID(server string) string {
	return fmt.Sprintf("bench_%d_%s", time.Now().UTC().UnixNano(), strings.ReplaceAll(server, ".", "_"))
}
