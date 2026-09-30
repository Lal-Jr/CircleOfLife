// Package application holds the Identity context's use cases: signing up, logging in and
// reading your own profile. Both successful sign-up and log-in return a session token.
package application

import (
	"context"
	"errors"
	"time"

	"circleoflife/internal/identity/domain"
)

// TokenIssuer issues the bearer token a member presents on later requests.
type TokenIssuer interface {
	Issue(id domain.UserID) (string, error)
}

// Profile is what a member sees about themselves.
type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type Service struct {
	users  domain.Users
	hasher domain.PasswordHasher
	tokens TokenIssuer
	now    func() time.Time
}

func NewService(users domain.Users, hasher domain.PasswordHasher, tokens TokenIssuer) *Service {
	return &Service{users: users, hasher: hasher, tokens: tokens, now: time.Now}
}

func (s *Service) SignUp(ctx context.Context, name, email, password string) (string, error) {
	addr, err := domain.NewEmail(email)
	if err != nil {
		return "", err
	}
	user, err := domain.Register(name, addr, password, s.hasher, s.now())
	if err != nil {
		return "", err
	}
	if err := s.users.Add(ctx, user); err != nil {
		return "", err
	}
	return s.tokens.Issue(user.ID())
}

// LogIn gives the same answer for an unknown email as for a wrong password, so it can't be
// used to find out who has an account.
func (s *Service) LogIn(ctx context.Context, email, password string) (string, error) {
	addr, err := domain.NewEmail(email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}
	user, err := s.users.ByEmail(ctx, addr)
	if errors.Is(err, domain.ErrUserNotFound) {
		return "", domain.ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}
	if err := user.Authenticate(password, s.hasher); err != nil {
		return "", err
	}
	return s.tokens.Issue(user.ID())
}

func (s *Service) Profile(ctx context.Context, id string) (*Profile, error) {
	user, err := s.users.ByID(ctx, domain.UserID(id))
	if err != nil {
		return nil, err
	}
	return &Profile{ID: string(user.ID()), Name: user.Name(), Email: string(user.Email()), CreatedAt: user.CreatedAt()}, nil
}
