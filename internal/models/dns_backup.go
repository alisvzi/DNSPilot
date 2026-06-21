package models

type DNSBackup struct {
	AdapterName string   `json:"adapter_name"`
	Servers     []string `json:"servers"`
	IsIPv6      bool     `json:"is_ipv6"`
	CreatedAt   string   `json:"created_at"`
}
