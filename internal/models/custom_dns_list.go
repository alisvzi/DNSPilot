package models

type CustomDNSList struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Servers     []DNSServer `json:"servers"`
}
