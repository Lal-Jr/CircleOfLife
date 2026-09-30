package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"circleoflife/internal/community/domain"
)

// In-memory fakes for the ports, so the use cases run without Postgres or Redis.

type fakePosts struct {
	posts    map[domain.PostID]*domain.Post
	comments []domain.Comment
	votes    map[domain.PostID]map[domain.MemberID]bool
}

func newFakePosts() *fakePosts {
	return &fakePosts{posts: map[domain.PostID]*domain.Post{}, votes: map[domain.PostID]map[domain.MemberID]bool{}}
}

func (f *fakePosts) Get(_ context.Context, id domain.PostID) (*domain.Post, error) {
	if p, ok := f.posts[id]; ok {
		return p, nil
	}
	return nil, domain.ErrPostNotFound
}
func (f *fakePosts) Add(_ context.Context, p *domain.Post) error { f.posts[p.ID()] = p; return nil }
func (f *fakePosts) AddComment(_ context.Context, c domain.Comment) error {
	f.comments = append(f.comments, c)
	return nil
}
func (f *fakePosts) ToggleHelpfulVote(_ context.Context, id domain.PostID, voter domain.MemberID) (bool, int, error) {
	if f.votes[id] == nil {
		f.votes[id] = map[domain.MemberID]bool{}
	}
	f.votes[id][voter] = !f.votes[id][voter]
	total := 0
	for _, v := range f.votes[id] {
		if v {
			total++
		}
	}
	return f.votes[id][voter], total, nil
}

type fakeFeed struct {
	FeedReader
	nearby []PostView
}

func (f *fakeFeed) Nearby(context.Context, FeedRequest) ([]PostView, error) { return f.nearby, nil }
func (f *fakeFeed) Comment(_ context.Context, id domain.CommentID) (*CommentView, error) {
	return &CommentView{ID: string(id)}, nil
}

type recordedEvents []domain.Event

func (r *recordedEvents) Publish(_ context.Context, events []domain.Event) {
	*r = append(*r, events...)
}

// stripTags stands in for the HTML sanitiser.
type stripTags struct{}

func (stripTags) Clean(s string) string { return strings.NewReplacer("<b>", "", "</b>", "").Replace(s) }

func newService() (*Service, *fakePosts, *fakeFeed, *recordedEvents) {
	posts, feed, events := newFakePosts(), &fakeFeed{}, &recordedEvents{}
	return NewService(posts, feed, events, stripTags{}), posts, feed, events
}

func TestPublishPostStoresCleanTextAndAnnouncesIt(t *testing.T) {
	svc, posts, _, events := newService()
	view, err := svc.PublishPost(context.Background(), PublishPost{
		Author: "m1", Title: "Need a <b>ladder</b>", Description: "Borrowing one for an hour", Kind: "help", Lat: 12.97, Lng: 77.59,
	})
	if err != nil {
		t.Fatal(err)
	}
	stored := posts.posts[domain.PostID(view.ID)]
	if stored == nil || stored.Title() != "Need a ladder" {
		t.Fatalf("stored post = %+v", stored)
	}
	if view.Title != "Need a ladder" || view.Lat != 12.97 || view.Lng != 77.59 || view.Type != "help" {
		t.Errorf("view = %+v", view)
	}
	if len(*events) != 1 {
		t.Fatalf("published %d events, want 1", len(*events))
	}
	if e, ok := (*events)[0].(domain.PostPublished); !ok || string(e.Post) != view.ID {
		t.Errorf("event = %#v", (*events)[0])
	}
}

func TestPublishPostRejectsTitleThatIsOnlyMarkup(t *testing.T) {
	svc, posts, _, events := newService()
	_, err := svc.PublishPost(context.Background(), PublishPost{Author: "m1", Title: "<b></b>", Description: "Something useful", Kind: "help", Lat: 1, Lng: 1})
	if !errors.Is(err, domain.ErrInvalidTitle) {
		t.Fatalf("err = %v, want ErrInvalidTitle", err)
	}
	if len(posts.posts) != 0 || len(*events) != 0 {
		t.Error("a rejected post was stored or announced")
	}
}

func TestCommentingOrVotingOnAMissingPostIsNotFound(t *testing.T) {
	svc, posts, _, _ := newService()
	ctx := context.Background()
	for _, id := range []string{"00000000-0000-0000-0000-000000000000", "not-a-uuid"} {
		if _, err := svc.AddComment(ctx, id, "m1", "hello"); !errors.Is(err, domain.ErrPostNotFound) {
			t.Errorf("AddComment(%q): err = %v", id, err)
		}
		if _, _, err := svc.ToggleHelpful(ctx, id, "m1"); !errors.Is(err, domain.ErrPostNotFound) {
			t.Errorf("ToggleHelpful(%q): err = %v", id, err)
		}
	}
	if len(posts.comments) != 0 || len(posts.votes) != 0 {
		t.Error("something was stored for a missing post")
	}
}

func TestToggleHelpfulAddsThenTakesBackAVote(t *testing.T) {
	svc, _, _, _ := newService()
	ctx := context.Background()
	view, _ := svc.PublishPost(ctx, PublishPost{Author: "m1", Title: "Walk group", Description: "Around the lake at six", Kind: "meetup", Lat: 1, Lng: 1})
	if voted, total, _ := svc.ToggleHelpful(ctx, view.ID, "m2"); !voted || total != 1 {
		t.Errorf("first toggle = %v, %d", voted, total)
	}
	if voted, total, _ := svc.ToggleHelpful(ctx, view.ID, "m2"); voted || total != 0 {
		t.Errorf("second toggle = %v, %d", voted, total)
	}
}

func TestFeedReportsWhetherAnotherPageFollows(t *testing.T) {
	svc, _, feed, _ := newService()
	feed.nearby = []PostView{{ID: "a"}, {ID: "b"}, {ID: "c"}}
	views, hasNext, err := svc.Feed(context.Background(), 1, 1, 5, 1, 2, "m1")
	if err != nil || !hasNext || len(views) != 2 {
		t.Errorf("limit 2 of 3: views=%d hasNext=%v err=%v", len(views), hasNext, err)
	}
	feed.nearby = feed.nearby[:2]
	views, hasNext, _ = svc.Feed(context.Background(), 1, 1, 5, 1, 2, "m1")
	if hasNext || len(views) != 2 {
		t.Errorf("limit 2 of exactly 2: views=%d hasNext=%v", len(views), hasNext)
	}
}
