package keyboard

import (
	"log"
	"strings"
	"sync"

	"LibreRemotePlay/src/devices"
	"github.com/go-vgo/robotgo"
	"github.com/pion/webrtc/v4"
)

var KeyboardEnabled = new(devices.DeviceEnabled).Disable()

func HandleKeyboard(d *webrtc.DataChannel) {

	if d.Label() != "keyboard" {
		return
	}

	d.OnOpen(func() {
		
		robotgo.KeySleep = 100
		
		log.Println("keyboard data channel is open")
		
	})

	keyState := make(map[string]bool)
	keyStateMutex := new(sync.RWMutex)

	d.OnMessage(func(msg webrtc.DataChannelMessage) {

		if !KeyboardEnabled.IsEnabled() {
			return
		}

		if !msg.IsString || msg.Data == nil {
			return
		}

		keyParts := strings.Split(string(msg.Data), "_")

		if len(keyParts) < 2 {
			return
		}

		key, keyExists := mapJSKeyToRobotGo(keyParts[0])

		if !keyExists {
			log.Println("keyboard key not found: ", keyParts[0])
			return
		}

		if keyParts[1] == "1" {
			keyStateMutex.RLock()
			if keyState[key] {
				keyStateMutex.RUnlock()
				return
			}
			keyStateMutex.RUnlock()
			keyStateMutex.Lock()
			keyState[key] = true
			keyStateMutex.Unlock()
			_ = robotgo.KeyDown(key)
			return
		} else {
			keyStateMutex.RLock()
			if !keyState[key] {
				keyStateMutex.RUnlock()
				return
			}
			keyStateMutex.RUnlock()
			keyStateMutex.Lock()
			keyState[key] = false
			keyStateMutex.Unlock()
			_ = robotgo.KeyUp(key)
			return
		}

	})

}
