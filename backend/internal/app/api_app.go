package app

import (
	"context"
	"net/http"

	platformconfig "github.com/your-org/agent-platform/internal/config"
	platformhttp "github.com/your-org/agent-platform/internal/http"
)

type APIApp struct {
	server *http.Server
}

func NewAPI(cfg platformconfig.Config) *APIApp {
	return &APIApp{
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: platformhttp.NewRouter(),
		},
	}
}

func (a *APIApp) Run(_ context.Context) error {
	return a.server.ListenAndServe()
}
