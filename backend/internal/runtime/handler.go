package runtime

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

func NewHandler(service Service, authorizer ...authz.Authorizer) *Handler {
	var selected authz.Authorizer
	if len(authorizer) > 0 {
		selected = authorizer[0]
	}
	return &Handler{service: service, authorizer: selected}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/projects/{projectId}/missions/{missionId}/agent-runtimes", h.listMissionRuntimes)
	r.Get("/api/projects/{projectId}/executor-sessions/{sessionId}", h.getSession)
	r.Get("/api/projects/{projectId}/executor-sessions/{sessionId}/events", h.listEvents)
	r.Get("/api/projects/{projectId}/executor-sessions/{sessionId}/transcript", h.getTranscript)
	r.Get("/api/projects/{projectId}/executor-sessions/{sessionId}/transcript/export", h.exportTranscript)
	r.Get("/api/projects/{projectId}/executor-sessions/{sessionId}/access-audits", h.listAccessAudits)
}

func (h *Handler) listMissionRuntimes(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ListMissionRuntimes(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) getSession(w http.ResponseWriter, r *http.Request) {
	session, err := h.service.GetSession(r.Context(), chi.URLParam(r, "sessionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *Handler) listEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.ListSessionEvents(r.Context(), chi.URLParam(r, "sessionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (h *Handler) getTranscript(w http.ResponseWriter, r *http.Request) {
	view := r.URL.Query().Get("view")
	if view == "" {
		view = "summary"
	}
	actorUserID := r.Header.Get("X-Actor-Id")
	if actorUserID == "" {
		actorUserID = "anonymous"
	}
	if view == "redacted" && h.authorizer != nil {
		if err := h.authorizer.Require(r.Context(), chi.URLParam(r, "projectId"), actorUserID, authz.CapabilityViewTranscripts); err != nil {
			http.Error(w, "redacted transcript access denied", http.StatusForbidden)
			return
		}
	}
	transcript, err := h.service.GetTranscriptView(r.Context(), chi.URLParam(r, "sessionId"), actorUserID, view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, transcript)
}

func (h *Handler) listAccessAudits(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ListAccessAudits(r.Context(), chi.URLParam(r, "sessionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) exportTranscript(w http.ResponseWriter, r *http.Request) {
	actorUserID := r.Header.Get("X-Actor-Id")
	if actorUserID == "" {
		actorUserID = "anonymous"
	}
	if h.authorizer != nil {
		if err := h.authorizer.Require(r.Context(), chi.URLParam(r, "projectId"), actorUserID, authz.CapabilityExportTranscript); err != nil {
			http.Error(w, "transcript export denied", http.StatusForbidden)
			return
		}
	}

	view, err := h.service.GetTranscriptView(r.Context(), chi.URLParam(r, "sessionId"), actorUserID, "redacted")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="transcript-export.json"`)
	_ = json.NewEncoder(w).Encode(view)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
