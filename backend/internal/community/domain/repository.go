package domain

import "context"

// Posts is the repository for the Post aggregate. It is defined here, in the domain, and
// implemented by the infrastructure layer.
type Posts interface {
	// Get loads a post, or returns ErrPostNotFound.
	Get(ctx context.Context, id PostID) (*Post, error)
	Add(ctx context.Context, p *Post) error
	AddComment(ctx context.Context, c Comment) error
	// ToggleHelpfulVote adds the member's vote, or takes it back if they had already voted,
	// and returns whether they now have a vote and the post's total.
	ToggleHelpfulVote(ctx context.Context, post PostID, voter MemberID) (voted bool, total int, err error)
}
