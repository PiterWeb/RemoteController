package bindings

import "net/http"

func openViGEmWizard() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
	
}