package document

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
	r.Route("/api/projects/{projectId}/missions/{missionId}/documents", func(r chi.Router) {
		r.Get("/", h.listDocuments)
		r.Post("/", h.createDocument)
		r.Get("/{documentId}/versions", h.listVersions)
		r.Post("/{documentId}/versions", h.createVersion)
		r.Post("/{documentId}/adopt", h.adoptVersion)
	})
}

type createDocumentRequest struct {
	Kind  string `json:"kind"`
	Title string `json:"title"`
}

type createVersionRequest struct {
	ParentVersionID   string `json:"parentVersionId"`
	ContentFormat     string `json:"contentFormat"`
	StorageKind       string `json:"storageKind"`
	ObjectKey         string `json:"objectKey"`
	ContentHash       string `json:"contentHash"`
	ContentText       string `json:"contentText"`
	ProducedByAgentID string `json:"producedByAgentId"`
	ProducedByRunID   string `json:"producedByRunId"`
	SourceRoundID     string `json:"sourceRoundId"`
}

type adoptVersionRequest struct {
	VersionID string `json:"versionId"`
}

// listDocuments 实现当前函数行为。
func (h *Handler) listDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := h.service.ListDocuments(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, documents)
}

// createDocument 实现当前函数行为。
func (h *Handler) createDocument(w http.ResponseWriter, r *http.Request) {
	if err := h.require(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req createDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	document, err := h.service.CreateDocument(r.Context(), CreateDocumentCmd{
		MissionID: chi.URLParam(r, "missionId"),
		Kind:      req.Kind,
		Title:     req.Title,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, document)
}

// listVersions 实现当前函数行为。
func (h *Handler) listVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := h.service.ListVersions(r.Context(), chi.URLParam(r, "documentId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, versions)
}

// createVersion 实现当前函数行为。
func (h *Handler) createVersion(w http.ResponseWriter, r *http.Request) {
	if err := h.require(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req createVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	version, err := h.service.CreateVersion(r.Context(), CreateVersionCmd{
		DocumentID:        chi.URLParam(r, "documentId"),
		ParentVersionID:   req.ParentVersionID,
		ContentFormat:     req.ContentFormat,
		StorageKind:       req.StorageKind,
		ObjectKey:         req.ObjectKey,
		ContentHash:       req.ContentHash,
		ContentText:       req.ContentText,
		ProducedByAgentID: req.ProducedByAgentID,
		ProducedByRunID:   req.ProducedByRunID,
		SourceRoundID:     req.SourceRoundID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, version)
}

// adoptVersion 实现当前函数行为。
func (h *Handler) adoptVersion(w http.ResponseWriter, r *http.Request) {
	if err := h.require(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var req adoptVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	version, err := h.service.AdoptVersion(r.Context(), chi.URLParam(r, "documentId"), req.VersionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, version)
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
