package gamepad

import (
	_ "embed"
	"fmt"
	"unsafe"
)

type clientVirtualGamepad uintptr

type EmulatedDevice struct {
	client clientVirtualGamepad
	pad    uintptr
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

func UpdateVirtualDevice(device EmulatedDevice, rg GamepadAPIXState, virtualState *ViGEmState) {

	// Get Real Input and convert to Virtual

	realState := gamepadAPIXToXInput(rg)

	realState.ToXInput(virtualState)

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
