package audio

import "LibreRemotePlay/src/devices"

type AudioProcess struct {
	Name string
	Pid  uint32 // It is not a real Pid on linux
}

var AudioEnabled = new(devices.DeviceEnabled).Enable()