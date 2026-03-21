package document

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

func (h *Handler) listDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := h.service.ListDocuments(r.Context(), chi.URLParam(r, "missionId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, documents)
}

func (h *Handler) createDocument(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) listVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := h.service.ListVersions(r.Context(), chi.URLParam(r, "documentId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, versions)
}

func (h *Handler) createVersion(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) adoptVersion(w http.ResponseWriter, r *http.Request) {
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

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
