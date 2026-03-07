package services

import (
	"context"

	"circleoflife/internal/models"
	"circleoflife/internal/repository"
)

type PostService interface {
	CreatePost(ctx context.Context, p *models.Post, lat, lng float64) error
	GetNearbyPosts(ctx context.Context, lat, lng float64, radius int) ([]models.Post, error)
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

func (s *postService) GetNearbyPosts(ctx context.Context, lat, lng float64, radius int) ([]models.Post, error) {
	return s.repo.GetNearbyPosts(ctx, lat, lng, radius)
}

func (s *postService) GetPostByID(ctx context.Context, id string, lat, lng float64) (*models.Post, error) {
	return s.repo.GetPostByID(ctx, id, lat, lng)
}
