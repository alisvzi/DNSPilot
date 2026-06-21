package models

type NetworkAdapter struct {
	ID string `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`

	MAC     string `json:"mac"`
	IPv4    string `json:"ipv4"`
	Gateway string `json:"gateway"`

	DNSServers []string `json:"dns_servers"`

	IsUp bool `json:"is_up"`
}
