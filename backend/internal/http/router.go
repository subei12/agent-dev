package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RouteRegistrar interface {
	RegisterRoutes(chi.Router)
}

func NewRouter(registrars ...RouteRegistrar) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", healthHandler())
	for _, registrar := range registrars {
		registrar.RegisterRoutes(r)
	}
	return r
}
