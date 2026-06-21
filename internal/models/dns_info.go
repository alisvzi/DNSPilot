package models

type DNSInfo struct {
	AdapterName string   `json:"adapter_name"`
	DNSServers  []string `json:"dns_servers"`
	IsActive    bool     `json:"is_active"`
}
