package dns

import "DNSPilot/internal/models"

func DefaultServers() []models.DNSServer {

	return []models.DNSServer{

		{
			ID:            "cloudflare",
			Name:          "Cloudflare",
			PrimaryIPv4:   "1.1.1.1",
			SecondaryIPv4: "1.0.0.1",
			Provider:      "Cloudflare",
			Description:   "Fast public DNS",
			Category:      models.PublicDNS,
			Enabled:       true,
		},

		{
			ID:            "google",
			Name:          "Google DNS",
			PrimaryIPv4:   "8.8.8.8",
			SecondaryIPv4: "8.8.4.4",
			Provider:      "Google",
			Description:   "Google Public DNS",
			Category:      models.PublicDNS,
			Enabled:       true,
		},

		{
			ID:            "quad9",
			Name:          "Quad9",
			PrimaryIPv4:   "9.9.9.9",
			SecondaryIPv4: "149.112.112.112",
			Provider:      "Quad9",
			Description:   "Secure DNS",
			Category:      models.SecureDNS,
			Enabled:       true,
		},
	}
}
