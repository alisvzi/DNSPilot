package models

type DNSCategory string

const (
	PublicDNS  DNSCategory = "public"
	PrivacyDNS DNSCategory = "privacy"
	SecureDNS  DNSCategory = "secure"
	FamilyDNS  DNSCategory = "family"
	CustomDNS  DNSCategory = "custom"
)
