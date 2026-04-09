package server

import (
	"net/http"
	"time"

	"github.com/tunakarakasoglu/finance-dashboard/config"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	registerRoutes(mux)
	return mux
}

func New(cfg config.Config) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

