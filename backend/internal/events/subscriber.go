package events

import (
	"context"
	"encoding/json"
	"log"

	"circleoflife/internal/cache"
	"circleoflife/internal/services"
)

func SubscribePostEvents(ctx context.Context) {
	if cache.Client == nil {
		log.Println("Redis client not initialized, skipping pub/sub listener")
		return
	}

	subscriber := cache.Client.Subscribe(ctx, PostEventsChannel)

	// Keep listening in a background goroutine
	go func() {
		defer subscriber.Close()
		log.Printf("Subscribed to Redis Channel: %s", PostEventsChannel)

		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				// Context cancelled or connection dropped
				log.Printf("Redis subscriber stopped: %v", err)
				return
			}

			var event services.EventPayload
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to unmarshal received event: %v", err)
				continue
			}

			// Broadcast received event physically originating from Redis down to connected SSE instances in memory
			log.Printf("Received Redis event, forwarding to SSE clients: %+v", event)
			services.Hub.Broadcast(event)
		}
	}()
}
