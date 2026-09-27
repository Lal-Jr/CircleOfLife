package repository

import (
	"context"

	"circleoflife/internal/db"
	"circleoflife/internal/models"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, c *models.Comment) error
	GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error)
}

type commentRepository struct{}

func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

func (r *commentRepository) CreateComment(ctx context.Context, c *models.Comment) error {
	// Joins in the author's name in the same round-trip, so the response to a
	// comment post matches the shape GetCommentsByPostID returns instead of
	// coming back with an empty authorName.
	query := `
		WITH inserted AS (
			INSERT INTO comments (post_id, user_id, content)
			VALUES ($1, $2, $3)
			RETURNING id, user_id, created_at
		)
		SELECT inserted.id, inserted.created_at, u.name
		FROM inserted
		JOIN users u ON u.id = inserted.user_id`

	err := db.Pool.QueryRow(ctx, query, c.PostID, c.UserID, c.Content).Scan(&c.ID, &c.CreatedAt, &c.AuthorName)
	return err
}

func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {
	query := `
		SELECT 
			c.id, c.post_id, c.user_id, c.content, c.created_at,
			u.name AS author_name 
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC`
		
	rows, err := db.Pool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.AuthorName); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
