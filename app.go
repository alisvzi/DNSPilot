package main

import (
	"context"

	"DNSPilot/internal/dns"
	"DNSPilot/internal/models"
	"DNSPilot/internal/services"
)

type App struct {
	ctx context.Context

	networkService *services.NetworkService
}

func NewApp() *App {
	return &App{
		networkService: services.NewNetworkService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetDefaultDNS() []models.DNSServer {
	return dns.DefaultServers()
}

func (a *App) GetNetworkAdapters() ([]models.NetworkAdapter, error) {
	return a.networkService.GetAdapters()
}
