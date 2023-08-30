package gamepad

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/PiterWeb/RemoteController/src/bin"
)

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

var (
	vigemDLL                      *syscall.LazyDLL
	vigem_free_proc               *syscall.LazyProc
	vigem_disconect_proc          *syscall.LazyProc
	vigem_alloc_proc              *syscall.LazyProc
	vigem_connect_proc            *syscall.LazyProc
	vigem_target_x360_alloc_proc  *syscall.LazyProc
	vigem_target_add_proc         *syscall.LazyProc
	vigem_target_x360_update_proc *syscall.LazyProc
	vigem_target_remove_proc      *syscall.LazyProc
	vigem_target_free_proc        *syscall.LazyProc
)

func init() {

	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	dllFile, err := os.Create("./ViGEmClient.dll")

	if err != nil {
		panic(err)
	}

	defer dllFile.Close()

	if strconv.IntSize == 32 {
		// x86 Architecture
		_, err = dllFile.Write(bin.ViGEmClient_x86)
	} else if strconv.IntSize == 64 {
		// x64 Architecture
		_, err = dllFile.Write(bin.ViGEmClient_x64)
	}

	if err != nil {
		panic(err)
	}

	exec.Command("regsvr32", path+"/gamepad/ViGEmClient.dll")

	vigemDLL = syscall.NewLazyDLL("ViGEmClient.dll")
	vigem_disconect_proc = vigemDLL.NewProc("vigem_disconnect")
	vigem_free_proc = vigemDLL.NewProc("vigem_free")
	vigem_alloc_proc = vigemDLL.NewProc("vigem_alloc")
	vigem_connect_proc = vigemDLL.NewProc("vigem_connect")
	vigem_target_x360_alloc_proc = vigemDLL.NewProc("vigem_target_x360_alloc")
	vigem_target_add_proc = vigemDLL.NewProc("vigem_target_add")
	vigem_target_remove_proc = vigemDLL.NewProc("vigem_target_remove")
	vigem_target_free_proc = vigemDLL.NewProc("vigem_target_free")
	vigem_target_x360_update_proc = vigemDLL.NewProc("vigem_target_x360_update")
}

func VIGEM_SUCCESS(val uintptr) bool {
	return val == uintptr(VIGEM_ERROR_NONE)
}

func handleVigemError(err error) error {

	if err != syscall.Errno(0) {
		return err
	}

	return nil

}

type ViGEmState struct {
	DwPacketNumber DWORD
	Gamepad        _ViGEm_GAMEPAD
}

type _ViGEm_GAMEPAD struct {
	wButtons      WORD
	bLeftTrigger  BYTE
	bRightTrigger BYTE
	sThumbLX      SHORT
	sThumbLY      SHORT
	sThumbRX      SHORT
	sThumbRY      SHORT
}

func (gamepad *_ViGEm_GAMEPAD) UpdateFromRawState(state RawControls) {

	gamepad.wButtons = WORD(state.Buttons)
	gamepad.bLeftTrigger = BYTE(state.LeftTrigger)
	gamepad.bRightTrigger = BYTE(state.RightTrigger)
	gamepad.sThumbLX = SHORT(state.ThumbLX)
	gamepad.sThumbLY = SHORT(state.ThumbLY)
	gamepad.sThumbRX = SHORT(state.ThumbRX)
	gamepad.sThumbRY = SHORT(state.ThumbRY)

}
