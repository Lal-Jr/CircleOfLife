package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"circleoflife/internal/cache"
	"circleoflife/internal/models"
	"circleoflife/internal/repository"
)

type PostService interface {
	CreatePost(ctx context.Context, p *models.Post, lat, lng float64) error
	GetNearbyPosts(ctx context.Context, lat, lng float64, radius, page, limit int) ([]models.Post, error)
	GetPostByID(ctx context.Context, id string, lat, lng float64) (*models.Post, error)
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) CreatePost(ctx context.Context, p *models.Post, lat, lng float64) error {
	return s.repo.CreatePost(ctx, p, lat, lng)
}

func (s *postService) GetNearbyPosts(ctx context.Context, lat, lng float64, radius, page, limit int) ([]models.Post, error) {
	// 1. Generate Cache Key using geographic buckets (rounded to 2 decimal places) and pagination boundaries
	latBucket := math.Round(lat*100) / 100
	lngBucket := math.Round(lng*100) / 100
	cacheKey := fmt.Sprintf("feed:%.2f:%.2f:%d:%d:%d", latBucket, lngBucket, radius, page, limit)

	// 2. Check for Cache Hit if Redis is available
	if cache.Client != nil {
		cachedData, err := cache.Client.Get(ctx, cacheKey).Result()
		if err == nil {
			var posts []models.Post
			if jsonErr := json.Unmarshal([]byte(cachedData), &posts); jsonErr == nil {
				log.Printf("feed cache hit (Key: %s)", cacheKey)
				return posts, nil
			} else {
				log.Printf("feed cache unmarshal failed: %v", jsonErr)
			}
		} else {
			// Intentionally ignoring redis.Nil error and letting it fallback gracefully
			log.Printf("feed cache miss (Key: %s)", cacheKey)
		}
	}

	// 3. Cache Miss - Cache Stampede Protection (Redis Lock)
	if cache.Client != nil {
		lockKey := "feed_lock:" + cacheKey
		// Attempt to acquire a short 5-second computation lock
		acquired, _ := cache.Client.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
		if !acquired {
			// Another request is actively building this cache.
			// To maintain raw performance, we gently yield and wait up to 500ms for it to finish.
			time.Sleep(200 * time.Millisecond)
			retryData, err := cache.Client.Get(ctx, cacheKey).Result()
			if err == nil {
				var retryPosts []models.Post
				if json.Unmarshal([]byte(retryData), &retryPosts) == nil {
					log.Printf("feed cache hit (Yielded to lock) (Key: %s)", cacheKey)
					return retryPosts, nil
				}
			}
			// If still empty after yielding, we proceed down to the DB intentionally to guarantee uptime
		} else {
			// Make sure we release our lock explicitly when we finish the DB fetch
			defer cache.Client.Del(ctx, lockKey)
		}
	}

	// 4. Fallback to Database Query natively via Repository PostGIS bounds
	posts, err := s.repo.GetNearbyPosts(ctx, lat, lng, radius, page, limit)
	if err != nil {
		return nil, err
	}

	// 4. Store the resulting model inside Redis synchronously handling any marshaling errors gracefully
	if cache.Client != nil && posts != nil {
		if jsonData, stringErr := json.Marshal(posts); stringErr == nil {
			// Expiration TTL set to 60 seconds internally mitigating load on hyper-localized feeds
			setErr := cache.Client.Set(ctx, cacheKey, jsonData, 60*time.Second).Err()
			if setErr != nil {
				log.Printf("Warning: failed to set feed cache: %v", setErr)
			} else {
				log.Printf("feed successfully cached (Key: %s)", cacheKey)
			}
		} else {
			log.Printf("Warning: failed to marshal posts for cache: %v", stringErr)
		}
	}

	return posts, nil
}

func (s *postService) GetPostByID(ctx context.Context, id string, lat, lng float64) (*models.Post, error) {
	return s.repo.GetPostByID(ctx, id, lat, lng)
}
