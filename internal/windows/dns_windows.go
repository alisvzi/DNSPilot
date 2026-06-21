//go:build windows

package windows

import (
	"net"
	"syscall"
	"unsafe"

	"DNSPilot/internal/models"

	"golang.org/x/sys/windows"
)

var (
	iphlpapi             = windows.NewLazySystemDLL("iphlpapi.dll")
	getAdaptersAddresses = iphlpapi.NewProc("GetAdaptersAddresses")
)

const (
	AF_UNSPEC               = 0
	GAA_FLAG_INCLUDE_PREFIX = 0x10
	OPER_STATUS_UP          = 1
)

type socketAddress struct {
	lpSockaddr      *syscall.RawSockaddrAny
	iSockaddrLength int32
}

type ipAdapterDnsServerAddress struct {
	Length   uint32
	Reserved uint32
	Next     *ipAdapterDnsServerAddress
	Address  socketAddress
}

type ipAdapterAddresses struct {
	Length      uint32
	IfIndex     uint32
	Next        *ipAdapterAddresses
	AdapterName *byte

	FirstUnicastAddress   unsafe.Pointer
	FirstAnycastAddress   unsafe.Pointer
	FirstMulticastAddress unsafe.Pointer
	FirstDnsServerAddress *ipAdapterDnsServerAddress

	DnsSuffix    *uint16
	Description  *uint16
	FriendlyName *uint16

	PhysicalAddress       [8]byte
	PhysicalAddressLength uint32
	Flags                 uint32
	Mtu                   uint32
	IfType                uint32
	OperStatus            uint32
	Ipv6IfIndex           uint32
	ZoneIndices           [16]uint32
	FirstPrefix           unsafe.Pointer
}

func GetSystemDNS() ([]models.DNSInfo, error) {
	var size uint32

	r1, _, _ := getAdaptersAddresses.Call(
		uintptr(AF_UNSPEC),
		uintptr(GAA_FLAG_INCLUDE_PREFIX),
		0,
		0,
		uintptr(unsafe.Pointer(&size)),
	)

	if size == 0 {
		if r1 == 0 {
			return nil, syscall.EINVAL
		}
		return nil, syscall.Errno(r1)
	}

	buffer := make([]byte, size)

	r1, _, _ = getAdaptersAddresses.Call(
		uintptr(AF_UNSPEC),
		uintptr(GAA_FLAG_INCLUDE_PREFIX),
		0,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if r1 != 0 {
		return nil, syscall.Errno(r1)
	}

	out := make([]models.DNSInfo, 0)

	for adapter := (*ipAdapterAddresses)(unsafe.Pointer(&buffer[0])); adapter != nil; adapter = adapter.Next {
		name := "Unknown Adapter"
		if adapter.FriendlyName != nil {
			name = windows.UTF16PtrToString(adapter.FriendlyName)
		} else if adapter.AdapterName != nil {
			name = windows.BytePtrToString(adapter.AdapterName)
		}

		out = append(out, models.DNSInfo{
			AdapterName: name,
			DNSServers:  parseDNSServers(adapter.FirstDnsServerAddress),
			IsActive:    adapter.OperStatus == OPER_STATUS_UP,
		})
	}

	return out, nil
}

func parseDNSServers(head *ipAdapterDnsServerAddress) []string {
	out := make([]string, 0)

	for cur := head; cur != nil; cur = cur.Next {
		if cur.Address.lpSockaddr == nil {
			continue
		}

		family := *(*uint16)(unsafe.Pointer(cur.Address.lpSockaddr))

		switch family {
		case 2: // AF_INET
			sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(cur.Address.lpSockaddr))
			ip := net.IPv4(sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3]).String()
			out = append(out, ip)

		case 23: // AF_INET6
			sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(cur.Address.lpSockaddr))
			ip := net.IP(sa.Addr[:]).String()
			out = append(out, ip)
		}
	}

	return out
}
