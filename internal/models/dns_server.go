package models

type DNSServer struct {
	ID string `json:"id"`

	Name string `json:"name"`

	PrimaryIPv4   string `json:"primary_ipv4"`
	SecondaryIPv4 string `json:"secondary_ipv4"`

	PrimaryIPv6   string `json:"primary_ipv6"`
	SecondaryIPv6 string `json:"secondary_ipv6"`

	Provider    string `json:"provider"`
	Description string `json:"description"`

	Category DNSCategory `json:"category"`

	Tags []string `json:"tags"`

	IsCustom bool `json:"is_custom"`
	Enabled  bool `json:"enabled"`
}
