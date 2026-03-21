package archive

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/your-org/agent-platform/internal/authz"
)

type Handler struct {
	service    Service
	authorizer authz.Authorizer
}

// NewHandler 创建并返回对应的组件。
func NewHandler(service Service, authorizer ...authz.Authorizer) *Handler {
	var selected authz.Authorizer
	if len(authorizer) > 0 {
		selected = authorizer[0]
	}
	return &Handler{service: service, authorizer: selected}
}

// RegisterRoutes 实现当前函数行为。
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/projects/{projectId}/missions/{missionId}/archive", h.createArchive)
	r.Get("/api/projects/{projectId}/missions/{missionId}/archive", h.getArchive)
}

// createArchive 实现当前函数行为。
func (h *Handler) createArchive(w http.ResponseWriter, r *http.Request) {
	if h.authorizer != nil {
		if err := h.authorizer.Require(r.Context(), chi.URLParam(r, "projectId"), r.Header.Get("X-Actor-Id"), authz.CapabilityManageArchive); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	}
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

// getArchive 实现当前函数行为。
func (h *Handler) getArchive(w http.ResponseWriter, r *http.Request) {
	archive, err := h.service.GetArchive(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, archive)
}

// writeJSON 实现当前函数行为。
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
