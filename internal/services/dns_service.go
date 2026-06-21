package services

import "DNSPilot/internal/windows"

type DNSService struct{}

func NewDNSService() *DNSService {
	return &DNSService{}
}

func (s *DNSService) GetSystemDNS() ([]string, error) {
	return windows.GetSystemDNS()
}
