package services

import (
	"context"
	"errors"

	"circleoflife/internal/models"
	"circleoflife/internal/repository"
	"circleoflife/pkg/auth"
)

type AuthService interface {
	Signup(ctx context.Context, req models.User) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Signup(ctx context.Context, req models.User) (string, error) {
	// 1. Hash password
	hashed, err := auth.HashPassword(req.PasswordHash) // Assuming the request struct puts plain text in PasswordHash temp natively before this point, or we parse it
	if err != nil {
		return "", err
	}
	req.PasswordHash = hashed

	// 2. Create User
	if err := s.userRepo.CreateUser(ctx, &req); err != nil {
		return "", errors.New("could not create user. email might already exist")
	}

	// 3. Generate Token
	token, err := auth.GenerateToken(req.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Fetch user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// 2. Compare password
	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid email or password")
	}

	// 3. Generate token
	token, err := auth.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
