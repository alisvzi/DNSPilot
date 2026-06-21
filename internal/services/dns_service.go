package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"DNSPilot/internal/models"
	"DNSPilot/internal/windows"
)

type DNSService struct {
	db *sql.DB
}

func NewDNSService(db *sql.DB) *DNSService {
	return &DNSService{db: db}
}

func (s *DNSService) GetSystemDNS() ([]models.DNSInfo, error) {
	return windows.GetSystemDNS()
}

func (s *DNSService) ApplyDNS(adapterName, primary, secondary string, ipv6 bool) error {
	if strings.TrimSpace(adapterName) == "" {
		return fmt.Errorf("adapter name is required")
	}
	if net.ParseIP(primary) == nil {
		return fmt.Errorf("invalid primary DNS: %s", primary)
	}
	if secondary != "" && net.ParseIP(secondary) == nil {
		return fmt.Errorf("invalid secondary DNS: %s", secondary)
	}

	family := "ipv4"
	if ipv6 {
		family = "ipv6"
	}

	setArgs := []string{
		"interface", family, "set", "dnsservers",
		fmt.Sprintf(`name="%s"`, adapterName),
		"static", primary, "primary",
	}
	if out, err := exec.Command("netsh", setArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("set dns failed: %w: %s", err, string(out))
	}

	if secondary != "" {
		addArgs := []string{
			"interface", family, "add", "dnsservers",
			fmt.Sprintf(`name="%s"`, adapterName),
			secondary, "index=2",
		}
		if out, err := exec.Command("netsh", addArgs...).CombinedOutput(); err != nil {
			return fmt.Errorf("add secondary dns failed: %w: %s", err, string(out))
		}
	}

	return nil
}

func (s *DNSService) FlushDNSCache() error {
	out, err := exec.Command("ipconfig", "/flushdns").CombinedOutput()
	if err != nil {
		return fmt.Errorf("flush dns cache failed: %w: %s", err, string(out))
	}
	return nil
}

func (s *DNSService) BackupCurrentDNS(adapterName string) error {
	if s.db == nil {
		return fmt.Errorf("database is not initialized")
	}
	if strings.TrimSpace(adapterName) == "" {
		return fmt.Errorf("adapter name is required")
	}

	all, err := s.GetSystemDNS()
	if err != nil {
		return err
	}

	var matched *models.DNSInfo
	for i := range all {
		if strings.EqualFold(all[i].AdapterName, adapterName) {
			matched = &all[i]
			break
		}
	}
	if matched == nil {
		return fmt.Errorf("adapter not found: %s", adapterName)
	}
	if len(matched.DNSServers) == 0 {
		return fmt.Errorf("adapter has no DNS servers configured")
	}

	serversJSON, err := json.Marshal(matched.DNSServers)
	if err != nil {
		return err
	}

	isIPv6 := false
	for _, s := range matched.DNSServers {
		if strings.Contains(s, ":") {
			isIPv6 = true
			break
		}
	}

	_, err = s.db.Exec(`
		INSERT INTO dns_backups (adapter_name, servers_json, is_ipv6, created_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(adapter_name) DO UPDATE SET
			servers_json = excluded.servers_json,
			is_ipv6 = excluded.is_ipv6,
			created_at = excluded.created_at
	`, adapterName, string(serversJSON), isIPv6, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}

func (s *DNSService) RestoreBackup(adapterName string) error {
	if s.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	row := s.db.QueryRow(`
		SELECT servers_json, is_ipv6
		FROM dns_backups
		WHERE adapter_name = ?
	`, adapterName)

	var serversJSON string
	var isIPv6 bool
	if err := row.Scan(&serversJSON, &isIPv6); err != nil {
		return err
	}

	var servers []string
	if err := json.Unmarshal([]byte(serversJSON), &servers); err != nil {
		return err
	}
	if len(servers) == 0 {
		return fmt.Errorf("backup is empty")
	}

	primary := servers[0]
	secondary := ""
	if len(servers) > 1 {
		secondary = servers[1]
	}

	return s.ApplyDNS(adapterName, primary, secondary, isIPv6)
}
