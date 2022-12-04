package server

import (
	"fmt"
	"net/http"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"time"
)

func Serve() {
	initNotion()
	initWx()

	svcConfig := config.GetConfig().Service
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", svcConfig.Host, svcConfig.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.SugaredLogger.Infof("Listening on %s", s.Addr)
	logger.SugaredLogger.Fatal(s.ListenAndServe())
}
