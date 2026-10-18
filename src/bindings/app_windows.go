package bindings

import (
	"LibreRemotePlay/src/devices/gamepad" // "LibreRemotePlay/src/devices/gamepad"
	"encoding/json/v2"
	"log"
	"net/http"
)

func openViGEmWizard() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		res := Res[string]{
			Result: gamepad.OpenViGEmWizard().Error(),
		}

		rawRes, err := json.Marshal(res)
		if err != nil {
			log.Println("Error mashalling res: ", err)
			w.WriteHeader(500)
			return
		}
		
		w.Write(rawRes)
	}
	
}