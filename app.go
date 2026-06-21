package main

import (
	"context"

	"DNSPilot/internal/dns"
	"DNSPilot/internal/models"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetDefaultDNS() []models.DNSServer {
	return dns.DefaultServers()
}
