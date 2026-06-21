package models

type NetworkAdapter struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Description string `json:"description"`

	DNSServers []string `json:"dns_servers"`

	IsUp bool `json:"is_up"`
}
