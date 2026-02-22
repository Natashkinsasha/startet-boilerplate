package server

import (
	"net/http"

	"starter-boilerplate/internal/shared/config"
)

func SetupMux() *http.ServeMux {
	return http.NewServeMux()
}

func SetupHTTPServer(mux *http.ServeMux, cfg config.AppConfig) *http.Server {
	return &http.Server{
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
