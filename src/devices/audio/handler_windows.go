package audio

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"LibreRemotePlay/src/bin"
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

func HandleAudio(ctx context.Context, track *webrtc.TrackLocalStaticSample) error {

	// Create exe temp file in temp dir
	StdoutPCMApplicationLoopbackFile, err := os.CreateTemp("", "StdoutPCMApplicationLoopback-*.exe")

	if err != nil {
		return err
	}

	defer os.Remove(StdoutPCMApplicationLoopbackFile.Name())

	if _, err := StdoutPCMApplicationLoopbackFile.Write(bin.StdoutPCMApplicationLoopback_exe); err != nil {
		return err
	}

	if err := StdoutPCMApplicationLoopbackFile.Close(); err != nil {
		return err
	}

	exePath := StdoutPCMApplicationLoopbackFile.Name()

	pid := ctx.Value("pid").(uint32)

	cmd := exec.CommandContext(ctx, exePath, fmt.Sprintf("%d", pid))

	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	buff := make([]byte, 1024)

	tickerChan := time.NewTicker(time.Millisecond * 20).C

	if err := cmd.Start(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tickerChan:
			if n, err := stdoutPipe.Read(buff); err != nil {
				return err
			} else if n == 0 {
				continue
			} else {
				log.Printf("Audio bytes readed: %d", n)
			}

			if AudioEnabled.IsEnabled() {
				if err = track.WriteSample(media.Sample{Data: buff, Duration: time.Millisecond * 20}); err != nil {
					return err
				}
			}
		}
	}

}

func GetAudioProcess() []AudioProcess {

	procs := []AudioProcess{
		{
			Name: "None",
			Pid:  0,
		},
	}

	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return procs
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return procs
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err := mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		return procs
	}

	var asm2 *wca.IAudioSessionManager2
	if err := mmd.Activate(wca.IID_IAudioSessionManager2, ole.CLSCTX_ALL, nil, &asm2); err != nil {
		return procs
	}

	defer asm2.Release()

	var ase *wca.IAudioSessionEnumerator
	if err := asm2.GetSessionEnumerator(&ase); err != nil {
		return procs
	}

	defer ase.Release()

	var sessionNumber int
	if err := ase.GetCount(&sessionNumber); err != nil {
		return procs
	}

	for i := range sessionNumber {
		var asc *wca.IAudioSessionControl
		if err := ase.GetSession(i, &asc); err != nil {
			continue
		}

		defer asc.Release()

		var asc2 *wca.IAudioSessionControl2
		if err := asc.PutQueryInterface(wca.IID_IAudioSessionControl2, &asc2); err != nil {
			continue
		}

		defer asc2.Release()

		var name string
		if err := asc.GetDisplayName(&name); err != nil {
			continue
		}

		var pid uint32 = 0
		if err := asc2.GetProcessId(&pid); err != nil {
			continue
		}

		log.Printf("Audio Src --- Name %s, pid: %d\n", name, pid)

		procs = append(procs, AudioProcess{
			Name: name,
			Pid:  pid,
		})

	}

	return procs
}
