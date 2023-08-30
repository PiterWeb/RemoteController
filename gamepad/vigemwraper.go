package gamepad

import "syscall"

type (
	BOOL          = uint32
	BOOLEAN       = byte
	BYTE          = byte
	DWORD         = uint32
	DWORD64       = uint64
	HANDLE        = uintptr
	HLOCAL        = uintptr
	LARGE_INTEGER = int64
	LONG          = int32
	LPVOID        = uintptr
	SIZE_T        = uintptr
	UINT          = uint32
	ULONG_PTR     = uintptr
	ULONGLONG     = uint64
	WORD          = uint16
	SHORT         = int16
)

type VIGEM_ERROR uintptr

const (
	VIGEM_ERROR_NONE VIGEM_ERROR = 0x20000000
)

func VIGEM_SUCCESS(val uintptr) bool {

	return val == uintptr(VIGEM_ERROR_NONE)
}

func handleVigemError(err error) error {

	if err != syscall.Errno(0) {
		return err
	}

	return nil

}

type _XInputState struct {
	dwPacketNumber DWORD
	Gamepad        _XINPUT_GAMEPAD
}

type XInputState _XInputState

type _XINPUT_GAMEPAD struct {
	wButtons      WORD
	bLeftTrigger  BYTE
	bRightTrigger BYTE
	sThumbLX      SHORT
	sThumbLY      SHORT
	sThumbRX      SHORT
	sThumbRY      SHORT
}

func (gamepad *_XINPUT_GAMEPAD) UpdateFromRawState(state StateRaw) {

	gamepad.wButtons = WORD(state.Buttons)
	gamepad.bLeftTrigger = BYTE(state.LeftTrigger)
	gamepad.bRightTrigger = BYTE(state.RightTrigger)
	gamepad.sThumbLX = SHORT(state.ThumbLX)
	gamepad.sThumbLY = SHORT(state.ThumbLY)
	gamepad.sThumbRX = SHORT(state.ThumbRX)
	gamepad.sThumbRY = SHORT(state.ThumbRY)

}
