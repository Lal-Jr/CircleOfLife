// Package infrastructure implements the Identity context's ports: members in Postgres,
// passwords with bcrypt and session tokens as JWTs.
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"circleoflife/internal/identity/application"
	"circleoflife/internal/identity/domain"
	"circleoflife/pkg/auth"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users struct {
	db *pgxpool.Pool
}

func NewUsers(db *pgxpool.Pool) *Users { return &Users{db: db} }

var _ domain.Users = (*Users)(nil)

const uniqueViolation = "23505"

func (r *Users) Add(ctx context.Context, u *domain.User) error {
	_, err := r.db.Exec(ctx, `INSERT INTO users (id, name, email, password_hash, created_at) VALUES ($1, $2, $3, $4, $5)`,
		string(u.ID()), u.Name(), string(u.Email()), u.PasswordHash(), u.CreatedAt())
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return domain.ErrEmailTaken
	}
	if err != nil {
		return fmt.Errorf("saving user: %w", err)
	}
	return nil
}

func (r *Users) ByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	return r.one(ctx, `WHERE email = $1`, string(email))
}

func (r *Users) ByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	return r.one(ctx, `WHERE id = $1`, string(id))
}

func (r *Users) one(ctx context.Context, where string, arg string) (*domain.User, error) {
	var (
		id, name, email, hash string
		created               time.Time
	)
	err := r.db.QueryRow(ctx, `SELECT id, name, email, password_hash, created_at FROM users `+where, arg).
		Scan(&id, &name, &email, &hash, &created)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("loading user: %w", err)
	}
	return domain.RehydrateUser(domain.UserID(id), name, domain.Email(email), hash, created), nil
}

// Bcrypt hashes passwords with bcrypt.
type Bcrypt struct{}

var _ domain.PasswordHasher = Bcrypt{}

func (Bcrypt) Hash(password string) (string, error) { return auth.HashPassword(password) }
func (Bcrypt) Matches(password, hash string) bool   { return auth.CheckPasswordHash(password, hash) }

// JWTIssuer signs session tokens with the server's secret.
type JWTIssuer struct {
	secret string
}

func NewJWTIssuer(secret string) JWTIssuer { return JWTIssuer{secret: secret} }

var _ application.TokenIssuer = JWTIssuer{}

func (j JWTIssuer) Issue(id domain.UserID) (string, error) {
	return auth.GenerateToken(string(id), j.secret)
}
