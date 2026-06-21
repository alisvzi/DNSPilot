package services

import (
	"DNSPilot/internal/models"
	"DNSPilot/internal/windows"
)

type DNSService struct{}

func NewDNSService() *DNSService {
	return &DNSService{}
}

func (s *DNSService) GetSystemDNS() ([]models.DNSInfo, error) {
	return windows.GetSystemDNS()
}
