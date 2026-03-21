package task

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
	r.Get("/api/projects/{projectId}/missions/{missionId}/task-board", h.getBoard)
	r.Post("/api/projects/{projectId}/missions/{missionId}/task-board", h.createBoard)
	r.Post("/api/projects/{projectId}/missions/{missionId}/tasks/{taskId}/claim", h.claimTask)
	r.Post("/api/projects/{projectId}/missions/{missionId}/tasks/{taskId}/handoffs", h.createHandoff)
	r.Post("/api/projects/{projectId}/missions/{missionId}/tasks/{taskId}/review-checkpoints", h.createCheckpoint)
}

type createBoardRequest struct {
	Title string              `json:"title"`
	Items []CreateTaskItemCmd `json:"items"`
}

type claimTaskRequest struct {
	AgentID     string `json:"agentId"`
	ClaimReason string `json:"claimReason"`
}

type createHandoffRequest struct {
	ToAgentID                string          `json:"toAgentId"`
	ToAdminAgent             bool            `json:"toAdminAgent"`
	Summary                  string          `json:"summary"`
	FromAgentID              string          `json:"fromAgentId"`
	OutputDocumentVersionIDs json.RawMessage `json:"outputDocumentVersionIds"`
	OutputRepoCandidateIDs   json.RawMessage `json:"outputRepoCandidateIds"`
	OutputArtifactVersionIDs json.RawMessage `json:"outputArtifactVersionIds"`
	ValidationSummary        string          `json:"validationSummary"`
	RiskSummary              string          `json:"riskSummary"`
	RecommendedNextAction    string          `json:"recommendedNextAction"`
}

type createCheckpointRequest struct {
	Kind                     string          `json:"kind"`
	RequestedByAgentID       string          `json:"requestedByAgentId"`
	AssignedAgentID          string          `json:"assignedAgentId"`
	Summary                  string          `json:"summary"`
	LinkedDocumentVersionIDs json.RawMessage `json:"linkedDocumentVersionIds"`
	LinkedRepoCandidateIDs   json.RawMessage `json:"linkedRepoCandidateIds"`
}

func (h *Handler) getBoard(w http.ResponseWriter, r *http.Request) {
	board, err := h.service.GetBoard(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, board)
}

func (h *Handler) createBoard(w http.ResponseWriter, r *http.Request) {
	var req createBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	board, err := h.service.CreateBoard(r.Context(), CreateBoardCmd{
		MissionID: chi.URLParam(r, "missionId"),
		Title:     req.Title,
		Items:     req.Items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, board)
}

func (h *Handler) claimTask(w http.ResponseWriter, r *http.Request) {
	var req claimTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claim, err := h.service.Claim(r.Context(), ClaimTaskCmd{
		TaskItemID:  chi.URLParam(r, "taskId"),
		AgentID:     req.AgentID,
		ClaimReason: req.ClaimReason,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, claim)
}

func (h *Handler) createHandoff(w http.ResponseWriter, r *http.Request) {
	var req createHandoffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	handoff, err := h.service.CreateHandoff(r.Context(), CreateHandoffCmd{
		TaskItemID:               chi.URLParam(r, "taskId"),
		FromAgentID:              req.FromAgentID,
		ToAgentID:                req.ToAgentID,
		ToAdminAgent:             req.ToAdminAgent,
		Summary:                  req.Summary,
		OutputDocumentVersionIDs: req.OutputDocumentVersionIDs,
		OutputRepoCandidateIDs:   req.OutputRepoCandidateIDs,
		OutputArtifactVersionIDs: req.OutputArtifactVersionIDs,
		ValidationSummary:        req.ValidationSummary,
		RiskSummary:              req.RiskSummary,
		RecommendedNextAction:    req.RecommendedNextAction,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, handoff)
}

func (h *Handler) createCheckpoint(w http.ResponseWriter, r *http.Request) {
	var req createCheckpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	checkpoint, err := h.service.RequestCheckpoint(r.Context(), RequestCheckpointCmd{
		MissionID:                chi.URLParam(r, "missionId"),
		TaskItemID:               chi.URLParam(r, "taskId"),
		Kind:                     req.Kind,
		RequestedByAgentID:       req.RequestedByAgentID,
		AssignedAgentID:          req.AssignedAgentID,
		Summary:                  req.Summary,
		LinkedDocumentVersionIDs: req.LinkedDocumentVersionIDs,
		LinkedRepoCandidateIDs:   req.LinkedRepoCandidateIDs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, checkpoint)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
