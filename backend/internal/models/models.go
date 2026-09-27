package models

import "time"

type PaginatedMeta struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	HasNext bool `json:"hasNext"`
}

type APIResponse struct {
	Data interface{}  `json:"data"`
	Meta *PaginatedMeta `json:"meta,omitempty"`
}

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
	Lat          float64   `json:"lat"`
	Lng          float64   `json:"lng"`
	// No `omitempty`: 0 is a meaningful distance (the viewer standing right
	// on top of the post), not an absent value. `omitempty` on a float64
	// treats 0.0 as empty and drops the field entirely, which made the
	// frontend render "NaNkm away" for any post at the viewer's exact
	// location (undefined / 1000 = NaN).
	Distance     float64   `json:"distance"`
	CommentCount int       `json:"commentCount"`
	HelpfulCount int       `json:"helpfulCount"`
	LikedByMe    bool      `json:"likedByMe"`
	Priority     string    `json:"priority,omitempty"`
}

type LikeResponse struct {
	Liked        bool `json:"liked"`
	HelpfulCount int  `json:"helpfulCount"`
}

type CreatePostInput struct {
	Title       string    `json:"title" validate:"required,min=5,max=120"`
	Description string    `json:"description" validate:"required,min=10,max=1000"`
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
	Content string `json:"content" validate:"required,max=500"`
}
