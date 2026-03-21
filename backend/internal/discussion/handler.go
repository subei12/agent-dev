package discussion

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
	r.Route("/api/projects/{projectId}/missions/{missionId}/discussions", func(r chi.Router) {
		r.Get("/", h.listSessions)
		r.Post("/", h.createSession)
		r.Post("/{sessionId}/rounds", h.createRound)
	})
}

type createSessionRequest struct {
	Topic              string `json:"topic"`
	InitiatedByAgentID string `json:"initiatedByAgentId"`
}

type createRoundRequest struct {
	PromptSummary             string          `json:"promptSummary"`
	ParticipantAgentIDs       json.RawMessage `json:"participantAgentIds"`
	Summary                   string          `json:"summary"`
	OpenQuestions             json.RawMessage `json:"openQuestions"`
	Conflicts                 json.RawMessage `json:"conflicts"`
	Conclusion                json.RawMessage `json:"conclusion"`
	AdoptedDocumentVersionIDs json.RawMessage `json:"adoptedDocumentVersionIds"`
}

// listSessions 实现当前函数行为。
func (h *Handler) listSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.service.ListSessions(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, sessions)
}

// createSession 实现当前函数行为。
func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	if err := h.require(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := h.service.CreateSession(r.Context(), CreateSessionCmd{
		MissionID:          chi.URLParam(r, "missionId"),
		Topic:              req.Topic,
		InitiatedByAgentID: req.InitiatedByAgentID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, session)
}

// createRound 实现当前函数行为。
func (h *Handler) createRound(w http.ResponseWriter, r *http.Request) {
	if err := h.require(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req createRoundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	round, err := h.service.CreateRound(r.Context(), CreateRoundCmd{
		SessionID:                 chi.URLParam(r, "sessionId"),
		PromptSummary:             req.PromptSummary,
		ParticipantAgentIDs:       req.ParticipantAgentIDs,
		Summary:                   req.Summary,
		OpenQuestions:             req.OpenQuestions,
		Conflicts:                 req.Conflicts,
		Conclusion:                req.Conclusion,
		AdoptedDocumentVersionIDs: req.AdoptedDocumentVersionIDs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, round)
}

// writeJSON 实现当前函数行为。
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

// require 实现当前函数行为。
func (h *Handler) require(r *http.Request) error {
	if h.authorizer == nil {
		return nil
	}
	return h.authorizer.Require(r.Context(), chi.URLParam(r, "projectId"), r.Header.Get("X-Actor-Id"), authz.CapabilityManageMission)
}
