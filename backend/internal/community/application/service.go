package application

import (
	"context"
	"time"

	"circleoflife/internal/community/domain"
)

type Service struct {
	posts  domain.Posts
	feed   FeedReader
	events EventPublisher
	text   TextPolicy
	now    func() time.Time
}

func NewService(posts domain.Posts, feed FeedReader, events EventPublisher, text TextPolicy) *Service {
	return &Service{posts: posts, feed: feed, events: events, text: text, now: time.Now}
}

// ---- Commands -------------------------------------------------------------------------------

type PublishPost struct {
	Author      domain.MemberID
	Title       string
	Description string
	Kind        string
	Lat, Lng    float64
	MeetupTime  *time.Time
}

// PublishPost pins a new post and tells the neighbourhood about it.
func (s *Service) PublishPost(ctx context.Context, cmd PublishPost) (*PostView, error) {
	at, err := domain.NewLocation(cmd.Lat, cmd.Lng)
	if err != nil {
		return nil, err
	}
	kind, err := domain.ParseKind(cmd.Kind)
	if err != nil {
		return nil, err
	}
	post, err := domain.PublishPost(cmd.Author, s.text.Clean(cmd.Title), s.text.Clean(cmd.Description), kind, at, cmd.MeetupTime, s.now())
	if err != nil {
		return nil, err
	}
	if err := s.posts.Add(ctx, post); err != nil {
		return nil, err
	}
	s.events.Publish(ctx, post.PullEvents())

	// A brand-new post has no author name, distance or activity to report yet.
	return &PostView{
		ID:          string(post.ID()),
		UserID:      string(post.Author()),
		Title:       post.Title(),
		Description: post.Description(),
		Type:        string(post.Kind()),
		MeetupTime:  post.MeetupTime(),
		CreatedAt:   post.CreatedAt(),
		Lat:         at.Lat(),
		Lng:         at.Lng(),
	}, nil
}

// AddComment replies to a post on a member's behalf.
func (s *Service) AddComment(ctx context.Context, postID string, author domain.MemberID, content string) (*CommentView, error) {
	id, err := domain.ParsePostID(postID)
	if err != nil {
		return nil, err
	}
	post, err := s.posts.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	comment, err := post.AddComment(author, s.text.Clean(content), s.now())
	if err != nil {
		return nil, err
	}
	if err := s.posts.AddComment(ctx, comment); err != nil {
		return nil, err
	}
	return s.feed.Comment(ctx, comment.ID())
}

// ToggleHelpful marks a post helpful for a member, or takes the mark back.
func (s *Service) ToggleHelpful(ctx context.Context, postID string, voter domain.MemberID) (voted bool, total int, err error) {
	id, err := domain.ParsePostID(postID)
	if err != nil {
		return false, 0, err
	}
	if _, err := s.posts.Get(ctx, id); err != nil {
		return false, 0, err
	}
	return s.posts.ToggleHelpfulVote(ctx, id, voter)
}

// ---- Queries --------------------------------------------------------------------------------

// Feed returns one page of posts around a point and whether another page follows.
func (s *Service) Feed(ctx context.Context, lat, lng float64, radiusKm, page, limit int, viewer domain.MemberID) ([]PostView, bool, error) {
	views, err := s.feed.Nearby(ctx, FeedRequest{At: Point{lat, lng}, RadiusKm: radiusKm, Page: page, Limit: limit, Viewer: viewer})
	if err != nil {
		return nil, false, err
	}
	hasNext := len(views) > limit
	if hasNext {
		views = views[:limit]
	}
	return views, hasNext, nil
}

// Post returns one post as seen from a point, or domain.ErrPostNotFound.
func (s *Service) Post(ctx context.Context, postID string, lat, lng float64, viewer domain.MemberID) (*PostView, error) {
	id, err := domain.ParsePostID(postID)
	if err != nil {
		return nil, err
	}
	return s.feed.Post(ctx, id, Point{lat, lng}, viewer)
}

// Comments lists a post's comments, oldest first. An unknown post simply has none.
func (s *Service) Comments(ctx context.Context, postID string) ([]CommentView, error) {
	id, err := domain.ParsePostID(postID)
	if err != nil {
		return nil, err
	}
	return s.feed.Comments(ctx, id)
}
