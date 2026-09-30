package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"circleoflife/internal/community/application"

	"github.com/redis/go-redis/v9"
)

const (
	feedTTL      = 60 * time.Second
	feedLockTTL  = 5 * time.Second
	feedLockWait = 200 * time.Millisecond
)

// CachedFeed caches feed pages in Redis in front of another FeedReader. Nearby viewers share a
// cache entry: coordinates are bucketed to two decimals (about 1 km). With no Redis client it
// passes straight through.
type CachedFeed struct {
	application.FeedReader
	redis *redis.Client
}

func NewCachedFeed(inner application.FeedReader, client *redis.Client) *CachedFeed {
	return &CachedFeed{FeedReader: inner, redis: client}
}

func (f *CachedFeed) Nearby(ctx context.Context, req application.FeedRequest) ([]application.PostView, error) {
	if f.redis == nil {
		return f.FeedReader.Nearby(ctx, req)
	}
	// The viewer is part of the key because each view carries a per-viewer likedByMe flag.
	key := fmt.Sprintf("feed:%.2f:%.2f:%d:%d:%d:%s",
		math.Round(req.At.Lat*100)/100, math.Round(req.At.Lng*100)/100, req.RadiusKm, req.Page, req.Limit, req.Viewer)

	if views, ok := f.get(ctx, key); ok {
		log.Printf("feed cache hit (Key: %s)", key)
		return views, nil
	}
	log.Printf("feed cache miss (Key: %s)", key)

	// Stampede protection: one request rebuilds an expired page while the others briefly wait
	// for it, then fall through to the database rather than fail.
	lockKey := "feed_lock:" + key
	if acquired, _ := f.redis.SetNX(ctx, lockKey, "1", feedLockTTL).Result(); acquired {
		defer f.redis.Del(ctx, lockKey)
	} else {
		time.Sleep(feedLockWait)
		if views, ok := f.get(ctx, key); ok {
			log.Printf("feed cache hit (Yielded to lock) (Key: %s)", key)
			return views, nil
		}
	}

	views, err := f.FeedReader.Nearby(ctx, req)
	if err != nil {
		return nil, err
	}
	if views != nil {
		if data, err := json.Marshal(views); err != nil {
			log.Printf("Warning: failed to marshal posts for cache: %v", err)
		} else if err := f.redis.Set(ctx, key, data, feedTTL).Err(); err != nil {
			log.Printf("Warning: failed to set feed cache: %v", err)
		}
	}
	return views, nil
}

func (f *CachedFeed) get(ctx context.Context, key string) ([]application.PostView, bool) {
	data, err := f.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}
	var views []application.PostView
	if err := json.Unmarshal([]byte(data), &views); err != nil {
		log.Printf("feed cache unmarshal failed: %v", err)
		return nil, false
	}
	return views, true
}
