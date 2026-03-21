package sse

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/projects/{projectId}/events/stream", h.stream)
}

func (h *Handler) stream(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel")
	if channel == "" {
		http.Error(w, "channel is required", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	events := h.hub.Subscribe(channel)
	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case event := <-events:
			fmt.Fprintf(w, "data: %s\n\n", event.Data)
			flusher.Flush()
		}
	}
}
