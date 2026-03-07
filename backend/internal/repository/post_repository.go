package repository

import (
	"context"
	"fmt"

	"circleoflife/internal/db"
	"circleoflife/internal/models"
)

type PostRepository interface {
	CreatePost(ctx context.Context, p *models.Post, lat, lng float64) error
	GetNearbyPosts(ctx context.Context, lat, lng float64, radiusKm int) ([]models.Post, error)
	GetPostByID(ctx context.Context, id string, lat, lng float64) (*models.Post, error)
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

func (r *postRepository) GetNearbyPosts(ctx context.Context, lat, lng float64, radiusKm int) ([]models.Post, error) {
	// Radius in meters
	radiusMeters := radiusKm * 1000

	// We calculate distance using PostGIS and also perform a spatial index filter via ST_DWithin
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.description, p.type, p.meetup_time, p.created_at,
			u.name AS author,
			ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance,
			(SELECT count(*) FROM comments c WHERE c.post_id = p.id) AS comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE ST_DWithin(
			p.location, 
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
			$3
		)
		ORDER BY distance ASC, p.created_at DESC
		LIMIT 50;
	`

	rows, err := db.Pool.Query(ctx, query, lng, lat, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("error querying nearby posts: %v", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Title, &p.Description, &p.Type, &p.MeetupTime, &p.CreatedAt,
			&p.Author, &p.Distance, &p.CommentCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %v", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id string, lat, lng float64) (*models.Post, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.description, p.type, p.meetup_time, p.created_at,
			u.name AS author,
			ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance,
			(SELECT count(*) FROM comments c WHERE c.post_id = p.id) AS comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $3
	`

	var p models.Post
	err := db.Pool.QueryRow(ctx, query, lng, lat, id).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Description, &p.Type, &p.MeetupTime, &p.CreatedAt,
		&p.Author, &p.Distance, &p.CommentCount,
	)

	if err != nil {
		return nil, fmt.Errorf("error retrieving single post: %v", err)
	}

	return &p, nil
}
