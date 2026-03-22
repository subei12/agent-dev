package agentcfg

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

// NewHandler 创建并返回 Agent 配置页面所需的 HTTP 处理器。
func NewHandler(service Service, authorizer ...authz.Authorizer) *Handler {
	var selected authz.Authorizer
	if len(authorizer) > 0 {
		selected = authorizer[0]
	}
	return &Handler{service: service, authorizer: selected}
}

// RegisterRoutes 注册 Agent 配置相关接口。
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/projects/{projectId}/agents", h.list)
	r.Patch("/api/projects/{projectId}/agents/{agentId}", h.update)
}

type updateRequest struct {
	Name           string   `json:"name"`
	Enabled        bool     `json:"enabled"`
	ExecutorType   string   `json:"executorType"`
	Command        string   `json:"command"`
	Args           []string `json:"args"`
	TimeoutSec     int32    `json:"timeoutSec"`
	MaxConcurrency int32    `json:"maxConcurrency"`
	AllowRepoRead  bool     `json:"allowRepoRead"`
	AllowRepoWrite bool     `json:"allowRepoWrite"`
	AllowNetwork   bool     `json:"allowNetwork"`
}

// list 实现当前函数行为。
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), chi.URLParam(r, "projectId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// update 实现当前函数行为。
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	if err := h.require(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.service.Update(r.Context(), UpdateAgentConfigCmd{
		AgentID:        chi.URLParam(r, "agentId"),
		Name:           req.Name,
		Enabled:        req.Enabled,
		ExecutorType:   req.ExecutorType,
		Command:        req.Command,
		Args:           req.Args,
		TimeoutSec:     req.TimeoutSec,
		MaxConcurrency: req.MaxConcurrency,
		AllowRepoRead:  req.AllowRepoRead,
		AllowRepoWrite: req.AllowRepoWrite,
		AllowNetwork:   req.AllowNetwork,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) require(r *http.Request) error {
	if h.authorizer == nil {
		return nil
	}
	return h.authorizer.Require(r.Context(), chi.URLParam(r, "projectId"), r.Header.Get("X-Actor-Id"), authz.CapabilityManageMission)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
