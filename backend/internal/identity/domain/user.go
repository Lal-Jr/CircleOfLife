// Package domain is the Identity context: who a member is and how they prove it. Other
// contexts only ever see a member's ID.
package domain

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var (
	ErrEmailTaken         = errors.New("could not create user. email might already exist")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidName        = errors.New("name can't be empty")
	ErrWeakPassword       = errors.New("password must be at least 6 characters")
)

const MinPasswordLength = 6

type UserID string

// Email is a member's sign-in address, a value object.
type Email string

func NewEmail(s string) (Email, error) {
	at := strings.IndexByte(s, '@')
	if at <= 0 || at == len(s)-1 || strings.ContainsAny(s, " \t\n") {
		return "", ErrInvalidEmail
	}
	return Email(s), nil
}

// PasswordHasher is how the domain turns passwords into something safe to store and checks
// them later. It's implemented with bcrypt in the infrastructure layer.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Matches(password, hash string) bool
}

// User is the aggregate root of the Identity context.
type User struct {
	id           UserID
	name         string
	email        Email
	passwordHash string
	createdAt    time.Time
}

// Register creates a member. The password is hashed straight away and never kept.
func Register(name string, email Email, password string, hasher PasswordHasher, now time.Time) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidName
	}
	if utf8.RuneCountInString(password) < MinPasswordLength {
		return nil, ErrWeakPassword
	}
	hash, err := hasher.Hash(password)
	if err != nil {
		return nil, err
	}
	return &User{id: UserID(uuid.NewString()), name: name, email: email, passwordHash: hash, createdAt: now}, nil
}

// RehydrateUser rebuilds a stored member.
func RehydrateUser(id UserID, name string, email Email, passwordHash string, createdAt time.Time) *User {
	return &User{id: id, name: name, email: email, passwordHash: passwordHash, createdAt: createdAt}
}

// Authenticate checks a sign-in attempt against the stored hash.
func (u *User) Authenticate(password string, hasher PasswordHasher) error {
	if !hasher.Matches(password, u.passwordHash) {
		return ErrInvalidCredentials
	}
	return nil
}

func (u *User) ID() UserID           { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() Email         { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) CreatedAt() time.Time { return u.createdAt }

// Users is the repository for the User aggregate.
type Users interface {
	// Add stores a new member, or returns ErrEmailTaken.
	Add(ctx context.Context, u *User) error
	// ByEmail and ByID return ErrUserNotFound when there's no such member.
	ByEmail(ctx context.Context, email Email) (*User, error)
	ByID(ctx context.Context, id UserID) (*User, error)
}
