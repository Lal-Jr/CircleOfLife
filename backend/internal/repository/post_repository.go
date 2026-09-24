package repository

import (
	"context"
	"fmt"

	"circleoflife/internal/db"
	"circleoflife/internal/models"
)

type PostRepository interface {
	CreatePost(ctx context.Context, p *models.Post, lat, lng float64) error
	GetNearbyPosts(ctx context.Context, lat, lng float64, radiusKm, page, limit int, viewerID string) ([]models.Post, error)
	GetPostByID(ctx context.Context, id string, lat, lng float64, viewerID string) (*models.Post, error)
	ToggleLike(ctx context.Context, postID, userID string) (liked bool, helpfulCount int, err error)
}

type postRepository struct{}

func NewPostRepository() PostRepository {
	return &postRepository{}
}

func (r *postRepository) CreatePost(ctx context.Context, p *models.Post, lat, lng float64) error {
	query := `
		INSERT INTO posts (user_id, title, description, type, location, meetup_time)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			ST_SetSRID(ST_MakePoint($5, $6), 4326),
			$7
		) RETURNING id, created_at
	`

	// Ensure MeetupTime is handled as nullable
	err := db.Pool.QueryRow(ctx, query, p.UserID, p.Title, p.Description, p.Type, lng, lat, p.MeetupTime).Scan(&p.ID, &p.CreatedAt)
	return err
}

func (r *postRepository) GetNearbyPosts(ctx context.Context, lat, lng float64, radiusKm, page, limit int, viewerID string) ([]models.Post, error) {
	// Radius in meters
	radiusMeters := radiusKm * 1000
	offset := (page - 1) * limit

	// We calculate distance using PostGIS and also perform a spatial index filter via ST_DWithin
	query := `
		SELECT
			p.id, p.user_id, p.title, p.description, p.type, p.meetup_time, p.created_at,
			u.name AS author,
			ST_Y(p.location::geometry) AS lat,
			ST_X(p.location::geometry) AS lng,
			ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance,
			(SELECT count(*) FROM comments c WHERE c.post_id = p.id) AS comment_count,
			(SELECT count(*) FROM post_likes pl WHERE pl.post_id = p.id) AS helpful_count,
			EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = $6) AS liked_by_me,
			(
				CASE
					WHEN ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) < 500 THEN 'high'
					WHEN ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) < 2000 THEN 'medium'
					ELSE 'normal'
				END
			) AS priority,
			(
				(1.0 / (EXTRACT(EPOCH FROM (NOW() - p.created_at))/3600 + 1.0)) +
				(1.0 / (ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) / 1000 + 1.0)) +
				(CASE WHEN ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) < 500 THEN 0.5 ELSE 0.0 END) +
				(CASE WHEN ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) < 2000 AND ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) >= 500 THEN 0.2 ELSE 0.0 END) +
				(CASE WHEN p.type = 'help' THEN 0.2 ELSE 0.0 END)
			) AS score
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE ST_DWithin(
			p.location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
			$3
		)
		ORDER BY score DESC, distance ASC, p.created_at DESC
		LIMIT $4 OFFSET $5;
	`

	rows, err := db.Pool.Query(ctx, query, lng, lat, radiusMeters, limit, offset, viewerID)
	if err != nil {
		return nil, fmt.Errorf("error querying nearby posts: %v", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		// The trailing `score` column drives ORDER BY only and isn't part of the
		// API response, but pgx still requires a destination for every selected
		// column, so it's scanned into a throwaway variable.
		var score float64
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Title, &p.Description, &p.Type, &p.MeetupTime, &p.CreatedAt,
			&p.Author, &p.Lat, &p.Lng, &p.Distance, &p.CommentCount, &p.HelpfulCount, &p.LikedByMe, &p.Priority, &score,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %v", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id string, lat, lng float64, viewerID string) (*models.Post, error) {
	query := `
		SELECT
			p.id, p.user_id, p.title, p.description, p.type, p.meetup_time, p.created_at,
			u.name AS author,
			ST_Y(p.location::geometry) AS lat,
			ST_X(p.location::geometry) AS lng,
			ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance,
			(SELECT count(*) FROM comments c WHERE c.post_id = p.id) AS comment_count,
			(SELECT count(*) FROM post_likes pl WHERE pl.post_id = p.id) AS helpful_count,
			EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = $4) AS liked_by_me
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $3
	`

	var p models.Post
	err := db.Pool.QueryRow(ctx, query, lng, lat, id, viewerID).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Description, &p.Type, &p.MeetupTime, &p.CreatedAt,
		&p.Author, &p.Lat, &p.Lng, &p.Distance, &p.CommentCount, &p.HelpfulCount, &p.LikedByMe,
	)

	if err != nil {
		return nil, fmt.Errorf("error retrieving single post: %v", err)
	}

	return &p, nil
}

func (r *postRepository) ToggleLike(ctx context.Context, postID, userID string) (bool, int, error) {
	var exists bool
	err := db.Pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`,
		postID, userID,
	).Scan(&exists)
	if err != nil {
		return false, 0, fmt.Errorf("error checking like state: %v", err)
	}

	if exists {
		_, err = db.Pool.Exec(ctx, `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`, postID, userID)
	} else {
		_, err = db.Pool.Exec(ctx, `INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2)`, postID, userID)
	}
	if err != nil {
		return false, 0, fmt.Errorf("error toggling like: %v", err)
	}

	var count int
	if err := db.Pool.QueryRow(ctx, `SELECT count(*) FROM post_likes WHERE post_id = $1`, postID).Scan(&count); err != nil {
		return false, 0, fmt.Errorf("error counting likes: %v", err)
	}

	return !exists, count, nil
}
