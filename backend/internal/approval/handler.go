package approval

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

// NewHandler 创建并返回对应的组件。
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes 实现当前函数行为。
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/projects/{projectId}/approvals/{approvalId}", h.getApproval)
	r.Get("/api/projects/{projectId}/missions/{missionId}/approvals", h.listApprovalsByMission)
	r.Post("/api/projects/{projectId}/repo-candidates/{candidateId}/publish", h.publishCandidate)
}

type publishCandidateRequest struct {
	MissionID string          `json:"missionId"`
	Action    string          `json:"action"`
	Intent    json.RawMessage `json:"intent"`
	CreatedBy string          `json:"createdBy"`
}

// getApproval 实现当前函数行为。
func (h *Handler) getApproval(w http.ResponseWriter, r *http.Request) {
	approval, err := h.service.Get(r.Context(), chi.URLParam(r, "approvalId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, approval)
}

// listApprovalsByMission 实现当前函数行为。
func (h *Handler) listApprovalsByMission(w http.ResponseWriter, r *http.Request) {
	approvals, err := h.service.ListByMission(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, approvals)
}

// publishCandidate 实现当前函数行为。
func (h *Handler) publishCandidate(w http.ResponseWriter, r *http.Request) {
	var req publishCandidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Action == "" {
		req.Action = "repo_candidate_publish"
	}

	approval, err := h.service.Create(r.Context(), CreateApprovalCmd{
		MissionID:   req.MissionID,
		Action:      req.Action,
		SubjectType: "repo_candidate_publish",
		SubjectID:   chi.URLParam(r, "candidateId"),
		Intent:      req.Intent,
		CreatedBy:   req.CreatedBy,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, approval)
}

// writeJSON 实现当前函数行为。
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
