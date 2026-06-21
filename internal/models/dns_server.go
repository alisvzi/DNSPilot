package models

type DNSServer struct {
	ID string `json:"id"`

	Name string `json:"name"`

	PrimaryIPv4 string `json:"primary_ipv4"`

	SecondaryIPv4 string `json:"secondary_ipv4"`

	Provider string `json:"provider"`

	Description string `json:"description"`

	Category DNSCategory `json:"category"`

	Enabled bool `json:"enabled"`
}
