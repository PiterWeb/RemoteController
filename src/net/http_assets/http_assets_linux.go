package http_assets

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

func InitHTTPAssets(clientPort int, assets embed.FS) error {

	staticFS, err := fs.Sub(assets, "frontend/build")

	if err != nil {
		return err
	}

	http.Handle("GET /", FileMiddleware(staticFS, http.FileServer(http.FS(staticFS))))

	err = http.ListenAndServe(fmt.Sprintf(":%d", clientPort), nil)

	if err != nil {
		return err
	}

	return nil

}

// If .html of the route is available it loads the .html otherwise try to load the given path
func FileMiddleware(staticFS fs.FS, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pathSplited := strings.SplitN(r.URL.Path, "/", 2)

		if len(pathSplited) != 2 {
			next.ServeHTTP(w, r)
			return
		}

		data, err := fs.ReadFile(staticFS, pathSplited[1]+".html")

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(200)
		w.Write(data)
	})
}
