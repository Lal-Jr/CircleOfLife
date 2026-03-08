package services

import (
	"log"
	"sync"
)

type EventType string

const (
	EventPostCreated EventType = "post_created"
)

type EventPayload struct {
	Type   EventType `json:"type"`
	PostID string    `json:"postId"`
	Lat    float64   `json:"lat"`
	Lng    float64   `json:"lng"`
}

type EventHub struct {
	clients map[chan EventPayload]bool
	mu      sync.RWMutex
}

var Hub *EventHub

func init() {
	Hub = &EventHub{
		clients: make(map[chan EventPayload]bool),
	}
}

// AddClient registers a new SSE connection channel
func (h *EventHub) AddClient() chan EventPayload {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientChan := make(chan EventPayload, 100)
	h.clients[clientChan] = true

	log.Printf("SSE Client connected. Total: %d", len(h.clients))
	return clientChan
}

// RemoveClient unregisters an SSE connection channel
func (h *EventHub) RemoveClient(clientChan chan EventPayload) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[clientChan]; ok {
		delete(h.clients, clientChan)
		close(clientChan)
		log.Printf("SSE Client disconnected. Total: %d", len(h.clients))
	}
}

// Broadcast sends an event payload to all currently connected clients
func (h *EventHub) Broadcast(event EventPayload) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for clientChan := range h.clients {
		select {
		case clientChan <- event:
		default:
			// If sending fails (e.g. channel full), we just drop the event for this client
			// to avoid blocking the broadcaster.
			log.Printf("Dropped event for sluggish SSE client")
		}
	}
}
