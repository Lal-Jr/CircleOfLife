package domain

import "time"

type CommentID string

// Comment is an entity inside the Post aggregate: it has its own identity but only exists as a
// reply to a post, so it can only be made with Post.AddComment.
type Comment struct {
	id        CommentID
	post      PostID
	author    MemberID
	content   string
	createdAt time.Time
}

func (c Comment) ID() CommentID        { return c.id }
func (c Comment) Post() PostID         { return c.post }
func (c Comment) Author() MemberID     { return c.author }
func (c Comment) Content() string      { return c.content }
func (c Comment) CreatedAt() time.Time { return c.createdAt }
