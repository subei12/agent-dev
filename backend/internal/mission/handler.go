package mission

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
	r.Route("/api/projects/{projectId}/missions", func(r chi.Router) {
		r.Post("/", h.createMission)
		r.Get("/{missionId}", h.getMission)
		r.Post("/{missionId}/decisions", h.decideMission)
	})
}

type createMissionRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	SourceType    string `json:"sourceType"`
	TeamID        string `json:"teamId"`
	RepoBindingID string `json:"repoBindingId"`
	AdminAgentID  string `json:"adminAgentId"`
	CreatedBy     string `json:"createdBy"`
}

type decideMissionRequest struct {
	DecidedByAgentID string `json:"decidedByAgentId"`
	Decision         string `json:"decision"`
	Summary          string `json:"summary"`
}

func (h *Handler) createMission(w http.ResponseWriter, r *http.Request) {
	if err := h.require(r, authz.CapabilityManageMission); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req createMissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mission, err := h.service.CreateMission(r.Context(), CreateMissionCmd{
		ProjectID:     chi.URLParam(r, "projectId"),
		Title:         req.Title,
		Description:   req.Description,
		SourceType:    req.SourceType,
		TeamID:        req.TeamID,
		RepoBindingID: req.RepoBindingID,
		AdminAgentID:  req.AdminAgentID,
		CreatedBy:     req.CreatedBy,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, mission)
}

func (h *Handler) getMission(w http.ResponseWriter, r *http.Request) {
	mission, err := h.service.GetMission(r.Context(), chi.URLParam(r, "projectId"), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, mission)
}

func (h *Handler) decideMission(w http.ResponseWriter, r *http.Request) {
	if err := h.require(r, authz.CapabilityManageMission); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req decideMissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decision, err := h.service.Decide(r.Context(), DecideMissionCmd{
		MissionID:        chi.URLParam(r, "missionId"),
		DecidedByAgentID: req.DecidedByAgentID,
		Decision:         req.Decision,
		Summary:          req.Summary,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, decision)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (h *Handler) require(r *http.Request, capability string) error {
	if h.authorizer == nil {
		return nil
	}
	return h.authorizer.Require(r.Context(), chi.URLParam(r, "projectId"), r.Header.Get("X-Actor-Id"), capability)
}
