package server

import (
	"net/http"

	"starter-boilerplate/internal/shared/config"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func SetupMux() *http.ServeMux {
	return http.NewServeMux()
}

func SetupHTTPServer(mux *http.ServeMux, cfg config.AppConfig) *http.Server {
	return &http.Server{
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
