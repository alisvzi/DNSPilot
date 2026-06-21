package services

import (
	"encoding/json"
	"os/exec"
)

type DNSService struct{}

func NewDNSService() *DNSService {
	return &DNSService{}
}

type AdapterDNS struct {
	Name       string   `json:"name"`
	DNSServers []string `json:"dns_servers"`
}

func (s *DNSService) GetCurrentDNS() ([]AdapterDNS, error) {

	cmd := exec.Command(
		"powershell",
		"-Command",
		`Get-DnsClientServerAddress |
ConvertTo-Json`,
	)

	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var raw interface{}

	err = json.Unmarshal(output, &raw)

	if err != nil {
		return nil, err
	}

	return []AdapterDNS{}, nil
}
