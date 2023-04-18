package server

import (
	"encoding/json"
	"fmt"
	"html/template"
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

// RenderHtml render htl response
func renderHtml(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(data.(string)))
}

// RenderJson render json response
func renderJson(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	dataJsonBytes, err := json.Marshal(data)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		logger.SugaredLogger.Infow("Can not marshal", "msg", data)
		//_, _ = fmt.Fprintf(w, data.(string))
		t := template.New("data")
		if err = t.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(dataJsonBytes)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}
