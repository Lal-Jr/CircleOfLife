package services

import (
	"context"

	"circleoflife/internal/models"
	"circleoflife/internal/repository"
)

type CommentService interface {
	CreateComment(ctx context.Context, c *models.Comment) error
	GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) CreateComment(ctx context.Context, c *models.Comment) error {
	return s.repo.CreateComment(ctx, c)
}

func (s *commentService) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {
	return s.repo.GetCommentsByPostID(ctx, postID)
}
