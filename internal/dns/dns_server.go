package models

type DNSServer struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	PrimaryIP   string   `json:"primary_ip"`
	SecondaryIP string   `json:"secondary_ip"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}
