package bin

import (
	_ "embed"
)

//go:embed ViGEmClient_x64.dll
var ViGEmClient_x64 []byte

//go:embed ViGEmClient_x86.dll
var ViGEmClient_x86 []byte

//go:embed ViGEmBus_1.22.0_x64_x86_arm64.exe
var ViGEmBus_exe []byte
