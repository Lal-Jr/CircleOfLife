package repository

import (
	"context"

	"circleoflife/internal/db"
	"circleoflife/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) CreateUser(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (name, email, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
	
	err := db.Pool.QueryRow(ctx, query, u.Name, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt)
	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at 
		FROM users 
		WHERE email = $1`
		
	var user models.User
	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, 
		&user.Name, 
		&user.Email, 
		&user.PasswordHash, 
		&user.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}
