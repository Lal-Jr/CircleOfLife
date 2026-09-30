// Package realtime is the supporting context that pushes live updates to open browsers. Other
// contexts publish events to a Redis channel; every server instance subscribes to that channel
// and fans each event out to its own Server-Sent Events clients, so a post made through one
// instance reaches viewers connected to any of them.
package realtime

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

const channel = "post_events"

type EventType string

const PostCreated EventType = "post_created"

// Event is the message sent to browsers, as JSON in an SSE "message" event.
type Event struct {
	Type   EventType `json:"type"`
	PostID string    `json:"postId"`
	Lat    float64   `json:"lat"`
	Lng    float64   `json:"lng"`
}

// Bus publishes events to every instance and delivers the ones it receives to local clients.
type Bus struct {
	redis *redis.Client // nil when Redis isn't configured: live updates are then off
	hub   *hub
}

func NewBus(client *redis.Client) *Bus {
	return &Bus{redis: client, hub: newHub()}
}

func (b *Bus) Publish(ctx context.Context, e Event) {
	if b.redis == nil {
		log.Println("Redis client not initialized, skipping publish")
		return
	}
	payload, err := json.Marshal(e)
	if err != nil {
		log.Printf("Failed to marshal event payload: %v", err)
		return
	}
	if err := b.redis.Publish(ctx, channel, string(payload)).Err(); err != nil {
		log.Printf("Failed to publish event to Redis: %v", err)
		return
	}
	log.Printf("%s event published to Redis", e.Type)
}

// Listen forwards events from the Redis channel to this instance's clients until ctx ends.
func (b *Bus) Listen(ctx context.Context) {
	if b.redis == nil {
		log.Println("Redis client not initialized, skipping pub/sub listener")
		return
	}
	sub := b.redis.Subscribe(ctx, channel)
	go func() {
		defer sub.Close()
		log.Printf("Subscribed to Redis Channel: %s", channel)
		for {
			msg, err := sub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Redis subscriber stopped: %v", err)
				return
			}
			var e Event
			if err := json.Unmarshal([]byte(msg.Payload), &e); err != nil {
				log.Printf("Failed to unmarshal received event: %v", err)
				continue
			}
			b.hub.broadcast(e)
		}
	}()
}

// hub holds this instance's connected SSE clients.
type hub struct {
	mu      sync.RWMutex
	clients map[chan Event]struct{}
}

func newHub() *hub { return &hub{clients: make(map[chan Event]struct{})} }

func (h *hub) add() chan Event {
	h.mu.Lock()
	defer h.mu.Unlock()
	c := make(chan Event, 100)
	h.clients[c] = struct{}{}
	log.Printf("SSE Client connected. Total: %d", len(h.clients))
	return c
}

func (h *hub) remove(c chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c)
		log.Printf("SSE Client disconnected. Total: %d", len(h.clients))
	}
}

func (h *hub) broadcast(e Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		select {
		case c <- e:
		default:
			// A client that isn't keeping up misses this event rather than blocking everyone.
			log.Printf("Dropped event for sluggish SSE client")
		}
	}
}
