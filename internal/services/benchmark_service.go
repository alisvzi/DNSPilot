package services

import (
	"DNSPilot/internal/benchmark"
	"DNSPilot/internal/models"
	"sort"
	"strings"
)

type BenchmarkService struct {
	profiles map[string]models.BenchmarkProfile
}

func NewBenchmarkService() *BenchmarkService {
	return &BenchmarkService{
		profiles: defaultBenchmarkProfiles(),
	}
}

func (s *BenchmarkService) GetBenchmarkProfiles() []models.BenchmarkProfile {
	out := make([]models.BenchmarkProfile, 0, len(s.profiles))
	for _, p := range s.profiles {
		out = append(out, p)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})

	return out
}

func (s *BenchmarkService) RunBenchmark(profileID string, servers []string) ([]models.BenchmarkResult, error) {
	profile, ok := s.profiles[strings.ToLower(profileID)]
	if !ok {
		profile = s.profiles["general"]
	}

	return benchmark.Run(profile, servers)
}

func defaultBenchmarkProfiles() map[string]models.BenchmarkProfile {
	return map[string]models.BenchmarkProfile{
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
	}
}
