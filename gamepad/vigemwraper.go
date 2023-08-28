package gamepad

type (
	BOOL          uint32
	BOOLEAN       byte
	BYTE          byte
	DWORD         uint32
	DWORD64       uint64
	HANDLE        uintptr
	HLOCAL        uintptr
	LARGE_INTEGER int64
	LONG          int32
	LPVOID        uintptr
	SIZE_T        uintptr
	UINT          uint32
	ULONG_PTR     uintptr
	ULONGLONG     uint64
	WORD          uint16
   )

// const vigem_targets_max = math.MaxUint16

// type VIGEM_TARGET_TYPE int

// const (
// 	xbox360Wired VIGEM_TARGET_TYPE = iota
// 	empty_vigem_type
// 	dualshock4Wired
// )

// type _VIGEM_CLIENT struct {
// 	hBusDevice                             unsafe.Pointer
// 	hDS4OutputReportPickupThread           unsafe.Pointer
// 	hDS4OutputReportPickupThreadAbortEvent unsafe.Pointer
// 	pTargetsList                           [vigem_targets_max]*_VIGEM_TARGET
// }

// type PVIGEM_CLIENT = *_VIGEM_CLIENT

// type DS4_OUTPUT_BUFFER struct {
// 	Buffer [64]uint8 // equivalent to [64]byte
// }

// type _VIGEM_TARGET struct {
// 	Sized                                float64
// 	SerialNo                             float64
// 	VIGEM_TARGET_STATE                   any
// 	VendorId                             uint16
// 	ProductId                            uint16
// 	Type                                 VIGEM_TARGET_TYPE
// 	Notification                         unsafe.Pointer
// 	NotificationUserData                 unsafe.Pointer
// 	IsWaitReadyUnsupported               bool
// 	CancelNotificationThreadEvent        unsafe.Pointer
// 	Ds4CachedOutputReport                DS4_OUTPUT_BUFFER
// 	Ds4CachedOutputReportUpdateAvailable unsafe.Pointer
// 	Ds4CachedOutputReportUpdateLock      unsafe.Pointer
// 	IsDisposing                          bool
// }

type VIGEM_ERROR uintptr

const (
	VIGEM_ERROR_NONE VIGEM_ERROR = 0x20000000
)

func VIGEM_SUCCESS(val uintptr) bool {

	return val == uintptr(VIGEM_ERROR_NONE)
}