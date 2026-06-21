package services

import (
	"errors"
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

type BenchmarkService struct {
	profiles map[string]models.BenchmarkProfile
}

func NewBenchmarkService() *BenchmarkService {
	return &BenchmarkService{
		profiles: map[string]models.BenchmarkProfile{
			"general": {
				ID:          "general",
				Name:        "General",
				Description: "Balanced DNS benchmark profile",
				Domains:     []string{"google.com", "cloudflare.com", "github.com"},
				Attempts:    5,
				TimeoutMS:   1500,
			},
			"gaming": {
				ID:          "gaming",
				Name:        "Gaming",
				Description: "Low-latency profile for game-related domains",
				Domains:     []string{"riotgames.com", "steampowered.com", "epicgames.com", "battle.net"},
				Attempts:    5,
				TimeoutMS:   1500,
			},
			"streaming": {
				ID:          "streaming",
				Name:        "Streaming",
				Description: "Profile for streaming services",
				Domains:     []string{"youtube.com", "netflix.com", "twitch.tv", "spotify.com"},
				Attempts:    5,
				TimeoutMS:   1500,
			},
			"development": {
				ID:          "development",
				Name:        "Development",
				Description: "Profile for dev tools and repositories",
				Domains:     []string{"github.com", "npmjs.com", "docker.com", "golang.org"},
				Attempts:    5,
				TimeoutMS:   1500,
			},
			"privacy": {
				ID:          "privacy",
				Name:        "Privacy",
				Description: "Profile for privacy-related services",
				Domains:     []string{"mozilla.org", "eff.org", "cloudflare.com"},
				Attempts:    5,
				TimeoutMS:   1500,
			},
		},
	}
}

func (s *BenchmarkService) GetBenchmarkProfiles() []models.BenchmarkProfile {
	out := make([]models.BenchmarkProfile, 0, len(s.profiles))
	for _, p := range s.profiles {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (s *BenchmarkService) RunBenchmark(profileID string, servers []string) ([]models.BenchmarkResult, error) {
	if len(servers) == 0 {
		return nil, errors.New("no DNS servers provided")
	}

	profile, ok := s.profiles[strings.ToLower(profileID)]
	if !ok {
		profile = s.profiles["general"]
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
				results <- s.benchmarkServer(server, profile)
			}
		}()
	}

	for _, server := range servers {
		jobs <- server
	}
	close(jobs)

	wg.Wait()
	close(results)

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

func (s *BenchmarkService) benchmarkServer(server string, profile models.BenchmarkProfile) models.BenchmarkResult {
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
		ID:          newID("bench"),
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
