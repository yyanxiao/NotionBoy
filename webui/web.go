package webui

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"notionboy/internal/pkg/logger"
	"os"
)

var (
	//go:embed all:dist
	distFS embed.FS

	distDirFS fs.FS
)

func init() {
	var err error
	distDirFS, err = fs.Sub(distFS, "dist")
	if err != nil {
		logger.SugaredLogger.Panicw("failed to load webui", "error", err)
	}
}

func RegisterHandlers(mux *http.ServeMux) {
	env := os.Getenv("ENV")
	logger.SugaredLogger.Infow("env", "env", env)
	if env == "dev" {
		mux.HandleFunc("/web", ReverseProxy)
		mux.HandleFunc("/web/", ReverseProxy)
	} else {
		mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.FS(distDirFS))))
	}
}

func ReverseProxy(w http.ResponseWriter, r *http.Request) {
	remote, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = r.Header
		req.Host = remote.Host
		req.URL = r.URL
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}

	// logger.SugaredLogger.Debugw("proxy request", "url", r.URL.String())
	proxy.ServeHTTP(w, r)
}
