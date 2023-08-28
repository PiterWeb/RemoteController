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
	vigemDLL                     *syscall.LazyDLL
	alloc_proc                   *syscall.LazyProc
	connect_proc                 *syscall.LazyProc
	vigem_target_x360_alloc_proc *syscall.LazyProc
	disconect_proc               *syscall.LazyProc
	free_proc                    *syscall.LazyProc
)

type clientVirtualGamepad uintptr

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
	disconect_proc = vigemDLL.NewProc("vigem_disconnect")
	free_proc = vigemDLL.NewProc("vigem_free")
	alloc_proc = vigemDLL.NewProc("vigem_alloc")
	connect_proc = vigemDLL.NewProc("vigem_connect")
	vigem_target_x360_alloc_proc = vigemDLL.NewProc("vigem_target_x360_alloc")
}

func InitializeEmulatedDevice() (clientVirtualGamepad, error) {

	client, _, err := alloc_proc.Call()

	if err != syscall.Errno(0) {
		return 0, err
	}

	if unsafe.Pointer(&client) == nil {
		return 0, fmt.Errorf("not enough memory to do that")
	}

	retval, _, err := connect_proc.Call(client)

	if err != syscall.Errno(0) {
		return 0, err
	}

	if !VIGEM_SUCCESS(retval) {
		return 0, fmt.Errorf("vigem bus connection failed with error code: 0x%X", retval)
	}

	return clientVirtualGamepad(client), nil
}

func EmulateDevice(device clientVirtualGamepad) error {

	return nil

}
