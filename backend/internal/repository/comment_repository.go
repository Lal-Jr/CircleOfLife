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
	query := `
		INSERT INTO comments (post_id, user_id, content) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
	
	err := db.Pool.QueryRow(ctx, query, c.PostID, c.UserID, c.Content).Scan(&c.ID, &c.CreatedAt)
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
