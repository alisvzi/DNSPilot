package main

import (
	"context"
	"database/sql"
	"log"

	"DNSPilot/internal/dns"
	"DNSPilot/internal/models"
	"DNSPilot/internal/services"
	"DNSPilot/internal/storage"
)

type App struct {
	ctx context.Context
	db  *sql.DB

	networkService   *services.NetworkService
	dnsService       *services.DNSService
	benchmarkService *services.BenchmarkService
	settingsService  *services.SettingsService
	listService      *services.ListService
}

func NewApp() *App {
	db, err := storage.Open()
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}
	if err := storage.Migrate(db); err != nil {
		log.Fatalf("migrate db failed: %v", err)
	}

	return &App{
		db:               db,
		networkService:   services.NewNetworkService(),
		dnsService:       services.NewDNSService(db),
		benchmarkService: services.NewBenchmarkService(),
		settingsService:  services.NewSettingsService(db),
		listService:      services.NewListService(db),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
}

func (a *App) GetDefaultDNS() []models.DNSServer {
	return dns.DefaultServers()
}

func (a *App) GetNetworkAdapters() ([]models.NetworkAdapter, error) {
	return a.networkService.GetAdapters()
}

func (a *App) GetSystemDNS() ([]models.DNSInfo, error) {
	return a.dnsService.GetSystemDNS()
}

func (a *App) ApplyDNS(adapterName, primary, secondary string, ipv6 bool) error {
	return a.dnsService.ApplyDNS(adapterName, primary, secondary, ipv6)
}

func (a *App) FlushDNSCache() error {
	return a.dnsService.FlushDNSCache()
}

func (a *App) BackupCurrentDNS(adapterName string) error {
	return a.dnsService.BackupCurrentDNS(adapterName)
}

func (a *App) RestoreBackup(adapterName string) error {
	return a.dnsService.RestoreBackup(adapterName)
}

func (a *App) GetBenchmarkProfiles() []models.BenchmarkProfile {
	return a.benchmarkService.GetBenchmarkProfiles()
}

func (a *App) RunBenchmark(profileID string, servers []string) ([]models.BenchmarkResult, error) {
	return a.benchmarkService.RunBenchmark(profileID, servers)
}

func (a *App) GetSettings() (models.Settings, error) {
	return a.settingsService.GetSettings()
}

func (a *App) SaveSettings(settings models.Settings) error {
	return a.settingsService.SaveSettings(settings)
}

func (a *App) GetCustomDNSLists() ([]models.CustomDNSList, error) {
	return a.listService.GetCustomDNSLists()
}

func (a *App) CreateCustomDNSList(name, description string) (string, error) {
	return a.listService.CreateCustomDNSList(name, description)
}

func (a *App) SaveCustomDNSList(list models.CustomDNSList) error {
	return a.listService.SaveCustomDNSList(list)
}

func (a *App) AddCustomDNSServer(listID string, server models.DNSServer) error {
	return a.listService.AddCustomDNSServer(listID, server)
}

func (a *App) DeleteCustomDNSServer(serverID string) error {
	return a.listService.DeleteCustomDNSServer(serverID)
}

func (a *App) DeleteCustomDNSList(id string) error {
	return a.listService.DeleteCustomDNSList(id)
}
