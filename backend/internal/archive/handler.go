package archive

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/projects/{projectId}/missions/{missionId}/archive", h.createArchive)
	r.Get("/api/projects/{projectId}/missions/{missionId}/archive", h.getArchive)
}

func (h *Handler) createArchive(w http.ResponseWriter, r *http.Request) {
	archive, manifest, err := h.service.CreateArchive(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"archive":  archive,
		"manifest": manifest,
	})
}

func (h *Handler) getArchive(w http.ResponseWriter, r *http.Request) {
	archive, err := h.service.GetArchive(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, archive)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
