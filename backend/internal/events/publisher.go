package events

import (
	"context"
	"encoding/json"
	"log"

	"circleoflife/internal/cache"
	"circleoflife/internal/services"
)

const PostEventsChannel = "post_events"

func PublishPostCreated(postID string, lat, lng float64) {
	if cache.Client == nil {
		log.Println("Redis client not initialized, skipping publish")
		return
	}

	event := services.EventPayload{
		Type:   services.EventPostCreated,
		PostID: postID,
		Lat:    lat,
		Lng:    lng,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event payload: %v", err)
		return
	}

	err = cache.Client.Publish(context.Background(), PostEventsChannel, string(payload)).Err()
	if err != nil {
		log.Printf("Failed to publish event to Redis: %v", err)
	} else {
		log.Println("post_created event successfully published to Redis")
	}
}
