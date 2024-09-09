package disk

// #cgo LDFLAGS: -lcryptsetup
// #include <libcryptsetup.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

type CryptStatusInfo int

const (
	CryptInvalid  CryptStatusInfo = C.CRYPT_INVALID
	CryptInactive CryptStatusInfo = C.CRYPT_INACTIVE
	CryptActive   CryptStatusInfo = C.CRYPT_ACTIVE
	CryptBusy     CryptStatusInfo = C.CRYPT_BUSY
)

func IsEncrypted(deviceName string) bool {
	cDeviceName := C.CString(deviceName)
	defer C.free(unsafe.Pointer(cDeviceName))

	state := CryptStatusInfo(C.crypt_status(nil, cDeviceName))

	switch state {
	case CryptInvalid:
		fallthrough
	case CryptInactive:
		return false
	case CryptBusy:
		fallthrough
	case CryptActive:
		return true
	}

	return false
}
