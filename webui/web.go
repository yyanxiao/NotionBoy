package webui

import (
	"embed"
	"io/fs"
	"net/http"

	"notionboy/internal/pkg/logger"
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
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.FS(distDirFS))))
}
