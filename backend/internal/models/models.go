package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name" validate:"required"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type Post struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"userId" db:"user_id"`
	Title        string    `json:"title" db:"title" validate:"required"`
	Description  string    `json:"description" db:"description" validate:"required"`
	Type         string    `json:"type" db:"type" validate:"required,oneof=help meetup"`
	MeetupTime   *time.Time`json:"meetupTime,omitempty" db:"meetup_time"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	
	// Appended fields used in API responses
	Author       string    `json:"author,omitempty"`
	Distance     float64   `json:"distance,omitempty"`
	CommentCount int       `json:"commentCount"`
}

type CreatePostInput struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Type        string    `json:"type" validate:"required,oneof=help meetup"`
	Lat         float64   `json:"lat" validate:"required,latitude"`
	Lng         float64   `json:"lng" validate:"required,longitude"`
	MeetupTime  *time.Time`json:"meetupTime"`
}

type Comment struct {
	ID         string    `json:"id" db:"id"`
	PostID     string    `json:"postId" db:"post_id"`
	UserID     string    `json:"userId" db:"user_id"`
	Content    string    `json:"content" db:"content" validate:"required"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	
	AuthorName string    `json:"authorName,omitempty"`
}

type CreateCommentInput struct {
	Content string `json:"content" validate:"required"`
}
