package services

import (
	"net"

	"DNSPilot/internal/models"
)

type NetworkService struct{}

func NewNetworkService() *NetworkService {
	return &NetworkService{}
}

func (s *NetworkService) GetAdapters() ([]models.NetworkAdapter, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	adapters := make([]models.NetworkAdapter, 0, len(interfaces))

	for _, iface := range interfaces {
		var ipv4 string

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				ipv4 = ipnet.IP.String()
				break
			}
		}

		adapters = append(adapters, models.NetworkAdapter{
			ID:          iface.Name,
			Name:        iface.Name,
			Description: iface.Name,
			MAC:         iface.HardwareAddr.String(),
			IPv4:        ipv4,
			IsUp:        iface.Flags&net.FlagUp != 0,
		})
	}

	return adapters, nil
}
