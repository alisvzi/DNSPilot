//go:build windows

package windows

import (
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

type IPAdapterAddresses struct {
	Length      uint32
	IfIndex     uint32
	Next        *IPAdapterAddresses
	AdapterName *byte

	FirstDnsServerAddress *syscall.RawSockaddrAny
}

func GetSystemDNS() ([]string, error) {

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

	return []string{"raw-data-ok"}, nil
}
