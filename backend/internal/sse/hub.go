package sse

import "sync"

type Event struct {
	Channel string
	Data    []byte
}

type Hub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[string][]chan Event),
	}
}

func (h *Hub) Subscribe(channel string) <-chan Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan Event, 8)
	h.subscribers[channel] = append(h.subscribers[channel], ch)
	return ch
}

func (h *Hub) Publish(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, ch := range h.subscribers[event.Channel] {
		select {
		case ch <- event:
		default:
		}
	}
}

func (h *Hub) PublishToChannel(channel string, data []byte) {
	h.Publish(Event{
		Channel: channel,
		Data:    data,
	})
}
