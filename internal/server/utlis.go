package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notionboy/internal/pkg/logger"
)

// func getQueryParams(req *http.Request, key string) []string {
// 	return req.URL.Query()[key]
// }

func getQueryParam(req *http.Request, key string) string {
	q := req.URL.Query()
	return q.Get(key)
}

// RenderError render error
func renderError(w http.ResponseWriter, code int, msg string, err error) {
	if msg == "" {
		logger.SugaredLogger.Errorw(msg, "err", err)
		http.Error(w, err.Error(), code)
		return
	}
	if err == nil {
		logger.SugaredLogger.Errorw(msg)
		http.Error(w, msg, code)
		return
	}
	http.Error(w, fmt.Sprintf("%s, error: %v", msg, err), code)
}

// RenderSuccess render success
func renderSuccess(w http.ResponseWriter, data interface{}) {
	dataJsonBytes, err := json.Marshal(data)
	if err != nil {
		logger.SugaredLogger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(dataJsonBytes)
}
