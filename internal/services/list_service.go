package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"DNSPilot/internal/models"
)

type ListService struct {
	db *sql.DB
}

func NewListService(db *sql.DB) *ListService {
	return &ListService{db: db}
}

func (s *ListService) GetCustomDNSLists() ([]models.CustomDNSList, error) {
	rows, err := s.db.Query(`SELECT id, name, description FROM custom_dns_lists ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []models.CustomDNSList

	for rows.Next() {
		var list models.CustomDNSList
		if err := rows.Scan(&list.ID, &list.Name, &list.Description); err != nil {
			return nil, err
		}

		servers, err := s.getServersByListID(list.ID)
		if err != nil {
			return nil, err
		}
		list.Servers = servers

		lists = append(lists, list)
	}

	return lists, rows.Err()
}

func (s *ListService) CreateCustomDNSList(name, description string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("list name is required")
	}

	id := newID("list")

	_, err := s.db.Exec(`
		INSERT INTO custom_dns_lists (id, name, description)
		VALUES (?, ?, ?)
	`, id, name, description)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *ListService) SaveCustomDNSList(list models.CustomDNSList) error {
	if strings.TrimSpace(list.Name) == "" {
		return fmt.Errorf("list name is required")
	}
	if list.ID == "" {
		list.ID = newID("list")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec(`
		INSERT INTO custom_dns_lists (id, name, description)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			description = excluded.description
	`, list.ID, list.Name, list.Description)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM custom_dns_servers WHERE list_id = ?`, list.ID); err != nil {
		return err
	}

	for _, server := range list.Servers {
		if err := validateCustomDNSServer(server); err != nil {
			return err
		}
		if server.ID == "" {
			server.ID = newID("dns")
		}
		if server.Category == "" {
			server.Category = models.CustomDNS
		}

		tagsJSON, err := json.Marshal(server.Tags)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO custom_dns_servers (
				id, list_id, name, primary_ipv4, secondary_ipv4, primary_ipv6, secondary_ipv6,
				provider, description, category, tags_json, is_custom, enabled
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			server.ID, list.ID, server.Name, server.PrimaryIPv4, server.SecondaryIPv4,
			server.PrimaryIPv6, server.SecondaryIPv6, server.Provider, server.Description,
			string(server.Category), string(tagsJSON), server.IsCustom, server.Enabled,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *ListService) AddCustomDNSServer(listID string, server models.DNSServer) error {
	listID = strings.TrimSpace(listID)
	if listID == "" {
		return fmt.Errorf("list id is required")
	}

	if err := validateCustomDNSServer(server); err != nil {
		return err
	}

	exists, err := s.customListExists(listID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("custom dns list not found: %s", listID)
	}

	if server.ID == "" {
		server.ID = newID("dns")
	}
	if server.Category == "" {
		server.Category = models.CustomDNS
	}
	if !server.Enabled {
		server.Enabled = true
	}

	tagsJSON, err := json.Marshal(server.Tags)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO custom_dns_servers (
			id, list_id, name, primary_ipv4, secondary_ipv4, primary_ipv6, secondary_ipv6,
			provider, description, category, tags_json, is_custom, enabled
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			list_id = excluded.list_id,
			name = excluded.name,
			primary_ipv4 = excluded.primary_ipv4,
			secondary_ipv4 = excluded.secondary_ipv4,
			primary_ipv6 = excluded.primary_ipv6,
			secondary_ipv6 = excluded.secondary_ipv6,
			provider = excluded.provider,
			description = excluded.description,
			category = excluded.category,
			tags_json = excluded.tags_json,
			is_custom = excluded.is_custom,
			enabled = excluded.enabled
	`,
		server.ID, listID, server.Name, server.PrimaryIPv4, server.SecondaryIPv4,
		server.PrimaryIPv6, server.SecondaryIPv6, server.Provider, server.Description,
		string(server.Category), string(tagsJSON), true, server.Enabled,
	)

	return err
}

func (s *ListService) DeleteCustomDNSServer(serverID string) error {
	serverID = strings.TrimSpace(serverID)
	if serverID == "" {
		return fmt.Errorf("server id is required")
	}

	_, err := s.db.Exec(`DELETE FROM custom_dns_servers WHERE id = ?`, serverID)
	return err
}

func (s *ListService) DeleteCustomDNSList(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	_, err := s.db.Exec(`DELETE FROM custom_dns_lists WHERE id = ?`, id)
	return err
}

func (s *ListService) getServersByListID(listID string) ([]models.DNSServer, error) {
	rows, err := s.db.Query(`
		SELECT id, name, primary_ipv4, secondary_ipv4, primary_ipv6, secondary_ipv6,
		       provider, description, category, tags_json, is_custom, enabled
		FROM custom_dns_servers
		WHERE list_id = ?
		ORDER BY name ASC
	`, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []models.DNSServer

	for rows.Next() {
		var (
			srv      models.DNSServer
			category string
			tagsJSON string
		)

		if err := rows.Scan(
			&srv.ID, &srv.Name, &srv.PrimaryIPv4, &srv.SecondaryIPv4, &srv.PrimaryIPv6, &srv.SecondaryIPv6,
			&srv.Provider, &srv.Description, &category, &tagsJSON, &srv.IsCustom, &srv.Enabled,
		); err != nil {
			return nil, err
		}

		srv.Category = models.DNSCategory(category)

		if tagsJSON != "" {
			_ = json.Unmarshal([]byte(tagsJSON), &srv.Tags)
		}

		servers = append(servers, srv)
	}

	return servers, rows.Err()
}

func (s *ListService) customListExists(listID string) (bool, error) {
	row := s.db.QueryRow(`SELECT 1 FROM custom_dns_lists WHERE id = ? LIMIT 1`, listID)

	var one int
	if err := row.Scan(&one); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func validateCustomDNSServer(server models.DNSServer) error {
	server.Name = strings.TrimSpace(server.Name)
	if server.Name == "" {
		return fmt.Errorf("server name is required")
	}

	hasPrimary := strings.TrimSpace(server.PrimaryIPv4) != "" || strings.TrimSpace(server.PrimaryIPv6) != ""
	if !hasPrimary {
		return fmt.Errorf("at least one primary DNS address is required")
	}

	if p := strings.TrimSpace(server.PrimaryIPv4); p != "" && net.ParseIP(p) == nil {
		return fmt.Errorf("invalid primary ipv4: %s", p)
	}
	if s := strings.TrimSpace(server.SecondaryIPv4); s != "" && net.ParseIP(s) == nil {
		return fmt.Errorf("invalid secondary ipv4: %s", s)
	}
	if p := strings.TrimSpace(server.PrimaryIPv6); p != "" && net.ParseIP(p) == nil {
		return fmt.Errorf("invalid primary ipv6: %s", p)
	}
	if s := strings.TrimSpace(server.SecondaryIPv6); s != "" && net.ParseIP(s) == nil {
		return fmt.Errorf("invalid secondary ipv6: %s", s)
	}

	return nil
}
