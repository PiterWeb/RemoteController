package http_assets

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

func InitHTTPAssets(clientPort int, assets embed.FS) error {

	staticFS, err := fs.Sub(assets, "frontend/build")

	if err != nil {
		return err
	}

	http.Handle("/", http.FileServer(http.FS(staticFS)))

	err = http.ListenAndServe(fmt.Sprintf(":%d", clientPort), nil)

	if err != nil {
		return err
	}

	return nil

}
