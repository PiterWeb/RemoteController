package gamepad

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/PiterWeb/RemoteController/bin"
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

type clientVirtualGamepad uintptr

type EmulatedDevice struct {
	client clientVirtualGamepad
	pad    uintptr
}

func init() {

	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	dll, err := os.Create("./ViGEmClient.dll")

	if err != nil {
		panic(err)
	}

	defer dll.Close()

	_, err = dll.Write(bin.ViGEmClient_x64)

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

func initializeEmulatedDevice() (clientVirtualGamepad, error) {

	client, _, err := vigem_alloc_proc.Call()

	err = handleVigemError(err)

	if err != nil {
		return 0, err
	}

	if unsafe.Pointer(&client) == nil {
		return 0, fmt.Errorf("not enough memory to do that")
	}

	retval, _, err := vigem_connect_proc.Call(client)

	err = handleVigemError(err)

	if err != nil {
		return 0, err
	}

	if !VIGEM_SUCCESS(retval) {
		return 0, fmt.Errorf("vigem bus connection failed with error code: 0x%X", retval)
	}

	return clientVirtualGamepad(client), nil
}

func UpdateVirtualDevice(device EmulatedDevice, realState State, virtualState *XInputState) {

	// Get Real Input and convert to Virtual
	realState.toXInput(virtualState)

	// Update the virtual gamepad
	vigem_target_x360_update_proc.Call(uintptr(device.client), device.pad, uintptr(unsafe.Pointer(&virtualState.Gamepad)))

}

func GenerateVirtualDevice() (EmulatedDevice, error) {

	device := new(EmulatedDevice)

	client, err := initializeEmulatedDevice()

	if err != nil {
		return *device, err
	}

	pad, _, err := vigem_target_x360_alloc_proc.Call()

	err = handleVigemError(err)

	if err != nil {
		return *device, err
	}

	device.client = client
	device.pad = pad

	pir, _, err := vigem_target_add_proc.Call(uintptr(client), pad)

	err = handleVigemError(err)

	if err != nil {
		return *device, err
	}

	if !VIGEM_SUCCESS(pir) {
		return *device, fmt.Errorf("target plugin failed with error code: 0x%x", pir)
	}

	return *device, nil

}

func FreeTargetAndDisconnect(device EmulatedDevice) {

	vigem_target_remove_proc.Call(uintptr(device.client), device.pad)
	vigem_target_free_proc.Call(device.pad)

	vigem_disconect_proc.Call(uintptr(device.client))
	vigem_free_proc.Call(uintptr(device.client))

}
