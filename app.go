package main

import (
	"context"

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

func (a *App) GetSystemDNS() ([]string, error) {
	return a.dnsService.GetSystemDNS()
}
