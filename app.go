package main

import (
	"context"

	"DNSPilot/internal/models"
	"DNSPilot/internal/services"
)

type App struct {
	ctx context.Context

	dnsService *services.DNSService
}

func NewApp() *App {
	return &App{
		dnsService: services.NewDNSService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetSystemDNS() ([]models.DNSInfo, error) {
	return a.dnsService.GetSystemDNS()
}
