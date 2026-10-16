// HTTP server
package http

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"

	"LibreRemotePlay/src/cli"
	"LibreRemotePlay/src/state"
)

func InitHTTPAssets(app state.App, serverMux *http.ServeMux) error {

	config := cli.GetConfig()
	clientPort := config.GetHTTPPort()

	var addr string

	if config.GetNetworkVisible() {
		addr = fmt.Sprintf("0.0.0.0:%d", clientPort)
	} else {
		addr = fmt.Sprintf("127.0.0.1:%d", clientPort)
	}

	staticFS, err := fs.Sub(app.Assets, "frontend/build")

	if err != nil {
		return err
	}
	
	serverMux.Handle("GET /", FileMiddleware(staticFS, http.FileServer(http.FS(staticFS))))

	httpServer := &http.Server{
		Handler: serverMux,
		Addr:    addr,
		BaseContext: func(l net.Listener) context.Context {
			return app.Ctx
		},
	}

	err = httpServer.ListenAndServe()

	if err != nil {
		errors.Join(err, errors.New("local web http server error"))
	}
	
	return err

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
		_, err = w.Write(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
