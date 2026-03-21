package discussion

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

func (h *Handler) listSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.service.ListSessions(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, sessions)
}

func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) createRound(w http.ResponseWriter, r *http.Request) {
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

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
