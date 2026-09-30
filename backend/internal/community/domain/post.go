package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// PostID identifies a post. Parsing rejects anything that isn't a UUID, so a malformed ID in a
// URL is simply a post that doesn't exist.
type PostID string

func ParsePostID(s string) (PostID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", ErrPostNotFound
	}
	return PostID(s), nil
}

// MemberID is a neighbour as the Community context sees them: just an identity issued by the
// Identity context. Names and emails belong to Identity and only appear here in read models.
type MemberID string

// Kind is what a post asks of the neighbourhood.
type Kind string

const (
	HelpRequest Kind = "help"   // someone needs a hand
	Meetup      Kind = "meetup" // an invitation to get together, usually at a set time
)

func ParseKind(s string) (Kind, error) {
	switch k := Kind(s); k {
	case HelpRequest, Meetup:
		return k, nil
	}
	return "", ErrInvalidKind
}

const (
	MaxTitleLength       = 120
	MaxDescriptionLength = 1000
	MaxCommentLength     = 500
)

// Post is the aggregate root of the Community context. Comments and helpful votes only exist
// on a post, so they are created through it; that is where the rules about them live.
type Post struct {
	id          PostID
	author      MemberID
	title       string
	description string
	kind        Kind
	location    Location
	meetupTime  *time.Time
	createdAt   time.Time

	events []Event
}

// PublishPost creates a new post and records a PostPublished event for the neighbourhood.
// Text is expected to be cleaned of markup already; a title or description that is empty after
// cleaning (say, one that was only tags) is rejected.
func PublishPost(author MemberID, title, description string, kind Kind, at Location, meetupTime *time.Time, now time.Time) (*Post, error) {
	if err := requireText(title, MaxTitleLength, ErrInvalidTitle, "title"); err != nil {
		return nil, err
	}
	if err := requireText(description, MaxDescriptionLength, ErrInvalidText, "description"); err != nil {
		return nil, err
	}
	if _, err := ParseKind(string(kind)); err != nil {
		return nil, err
	}
	p := &Post{
		id:          PostID(uuid.NewString()),
		author:      author,
		title:       title,
		description: description,
		kind:        kind,
		location:    at,
		meetupTime:  meetupTime,
		createdAt:   now,
	}
	p.events = append(p.events, PostPublished{Post: p.id, At: at})
	return p, nil
}

// RehydratePost rebuilds a stored post. Stored posts were valid when published, so no rules run.
func RehydratePost(id PostID, author MemberID, title, description string, kind Kind, at Location, meetupTime *time.Time, createdAt time.Time) *Post {
	return &Post{id: id, author: author, title: title, description: description, kind: kind, location: at, meetupTime: meetupTime, createdAt: createdAt}
}

// AddComment records a neighbour's reply to this post.
func (p *Post) AddComment(author MemberID, content string, now time.Time) (Comment, error) {
	if err := requireText(content, MaxCommentLength, ErrInvalidComment, "comment"); err != nil {
		return Comment{}, err
	}
	return Comment{id: CommentID(uuid.NewString()), post: p.id, author: author, content: content, createdAt: now}, nil
}

// PullEvents hands over the events recorded since the post was loaded or created, once.
func (p *Post) PullEvents() []Event {
	events := p.events
	p.events = nil
	return events
}

func (p *Post) ID() PostID             { return p.id }
func (p *Post) Author() MemberID       { return p.author }
func (p *Post) Title() string          { return p.title }
func (p *Post) Description() string    { return p.description }
func (p *Post) Kind() Kind             { return p.kind }
func (p *Post) Location() Location     { return p.location }
func (p *Post) MeetupTime() *time.Time { return p.meetupTime }
func (p *Post) CreatedAt() time.Time   { return p.createdAt }

func requireText(s string, max int, kind error, what string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("%w: %s can't be empty", kind, what)
	}
	if len([]rune(s)) > max {
		return fmt.Errorf("%w: %s is longer than %d characters", kind, what, max)
	}
	return nil
}
