// Package application holds the Community context's use cases. Commands change the Post
// aggregate through its repository; queries read projections built for the feed. Everything
// outside the process (database, cache, event bus, HTML sanitiser) is reached through the
// ports in this file.
package application

import (
	"context"
	"time"

	"circleoflife/internal/community/domain"
)

// PostView is the read model of a post as a neighbour sees it: the post plus its author's
// name, how far away it is, its activity counts and whether the viewer marked it helpful.
type PostView struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	MeetupTime  *time.Time `json:"meetupTime,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	Author      string     `json:"author,omitempty"`
	Lat         float64    `json:"lat"`
	Lng         float64    `json:"lng"`
	// No `omitempty`: 0 is a meaningful distance (the viewer standing right on top of the
	// post), not an absent value; dropping it made the frontend render "NaNkm away".
	Distance     float64 `json:"distance"`
	CommentCount int     `json:"commentCount"`
	HelpfulCount int     `json:"helpfulCount"`
	LikedByMe    bool    `json:"likedByMe"`
	Priority     string  `json:"priority,omitempty"`
}

// CommentView is the read model of a comment, with its author's name.
type CommentView struct {
	ID         string    `json:"id"`
	PostID     string    `json:"postId"`
	UserID     string    `json:"userId"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	AuthorName string    `json:"authorName,omitempty"`
}

// Point is where a viewer is looking from. Reads never store it, so unlike domain.Location it
// isn't validated: coordinates out of range simply match nothing.
type Point struct {
	Lat, Lng float64
}

// FeedRequest asks for one page of the posts around a point, ranked by domain.Score.
type FeedRequest struct {
	At       Point
	RadiusKm int
	Page     int
	Limit    int
	Viewer   domain.MemberID
}

// FeedReader serves the read side. Nearby returns up to Limit+1 views so the caller can tell
// whether another page exists.
type FeedReader interface {
	Nearby(ctx context.Context, req FeedRequest) ([]PostView, error)
	Post(ctx context.Context, id domain.PostID, from Point, viewer domain.MemberID) (*PostView, error)
	Comments(ctx context.Context, post domain.PostID) ([]CommentView, error)
	Comment(ctx context.Context, id domain.CommentID) (*CommentView, error)
}

// EventPublisher passes domain events on to the rest of the system.
type EventPublisher interface {
	Publish(ctx context.Context, events []domain.Event)
}

// TextPolicy strips markup from what neighbours type, so stored text is always plain.
type TextPolicy interface {
	Clean(s string) string
}
