package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

var now = time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)

func mustLocation(t *testing.T, lat, lng float64) Location {
	t.Helper()
	at, err := NewLocation(lat, lng)
	if err != nil {
		t.Fatal(err)
	}
	return at
}

func TestNewLocationRejectsOutOfRangeCoordinates(t *testing.T) {
	for _, c := range []struct{ lat, lng float64 }{{91, 0}, {-91, 0}, {0, 181}, {0, -181}} {
		if _, err := NewLocation(c.lat, c.lng); !errors.Is(err, ErrInvalidLocation) {
			t.Errorf("NewLocation(%v, %v) = %v, want ErrInvalidLocation", c.lat, c.lng, err)
		}
	}
	if _, err := NewLocation(12.97, 77.59); err != nil {
		t.Errorf("valid location rejected: %v", err)
	}
}

func TestPublishPostRecordsPostPublished(t *testing.T) {
	at := mustLocation(t, 12.97, 77.59)
	p, err := PublishPost("m1", "Need a ladder", "Borrowing a ladder for an hour", HelpRequest, at, nil, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParsePostID(string(p.ID())); err != nil {
		t.Errorf("post ID %q is not a UUID", p.ID())
	}
	events := p.PullEvents()
	if len(events) != 1 || events[0] != (PostPublished{Post: p.ID(), At: at}) {
		t.Fatalf("events = %#v, want one PostPublished", events)
	}
	if again := p.PullEvents(); len(again) != 0 {
		t.Errorf("events were handed over twice: %#v", again)
	}
}

func TestPublishPostEnforcesItsRules(t *testing.T) {
	at := mustLocation(t, 12.97, 77.59)
	cases := []struct {
		name               string
		title, description string
		kind               Kind
		want               error
	}{
		{"blank title", "   ", "Borrowing a ladder", HelpRequest, ErrInvalidTitle},
		{"title too long", strings.Repeat("a", MaxTitleLength+1), "Borrowing a ladder", HelpRequest, ErrInvalidTitle},
		{"blank description", "Need a ladder", "", HelpRequest, ErrInvalidText},
		{"description too long", "Need a ladder", strings.Repeat("é", MaxDescriptionLength+1), HelpRequest, ErrInvalidText},
		{"unknown kind", "Need a ladder", "Borrowing a ladder", Kind("sale"), ErrInvalidKind},
	}
	for _, c := range cases {
		if _, err := PublishPost("m1", c.title, c.description, c.kind, at, nil, now); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
	// Limits count characters, not bytes.
	if _, err := PublishPost("m1", strings.Repeat("é", MaxTitleLength), "Borrowing a ladder", Meetup, at, nil, now); err != nil {
		t.Errorf("title of exactly %d characters rejected: %v", MaxTitleLength, err)
	}
}

func TestAddCommentBelongsToThePost(t *testing.T) {
	p := RehydratePost("11111111-1111-1111-1111-111111111111", "author", "Title", "Description", Meetup, mustLocation(t, 1, 2), nil, now)
	c, err := p.AddComment("neighbour", "I'll be there", now)
	if err != nil {
		t.Fatal(err)
	}
	if c.Post() != p.ID() || c.Author() != "neighbour" || c.Content() != "I'll be there" || !c.CreatedAt().Equal(now) {
		t.Errorf("comment = %+v", c)
	}
	if _, err := p.AddComment("neighbour", " ", now); !errors.Is(err, ErrInvalidComment) {
		t.Errorf("blank comment: err = %v", err)
	}
	if _, err := p.AddComment("neighbour", strings.Repeat("x", MaxCommentLength+1), now); !errors.Is(err, ErrInvalidComment) {
		t.Errorf("long comment: err = %v", err)
	}
}

func TestParsePostIDTreatsMalformedIDsAsMissing(t *testing.T) {
	if _, err := ParsePostID("not-a-uuid"); !errors.Is(err, ErrPostNotFound) {
		t.Errorf("err = %v, want ErrPostNotFound", err)
	}
}

func TestProximityThresholds(t *testing.T) {
	for _, c := range []struct {
		meters float64
		want   Proximity
	}{{0, Urgent}, {499.9, Urgent}, {500, Nearby}, {1999.9, Nearby}, {2000, Distant}, {9000, Distant}} {
		if got := ProximityOf(c.meters); got != c.want {
			t.Errorf("ProximityOf(%v) = %v, want %v", c.meters, got, c.want)
		}
	}
}

func TestScoreRanksCloseRecentHelpRequestsFirst(t *testing.T) {
	hour := time.Hour
	if Score(hour, 300, Meetup) <= Score(hour, 3000, Meetup) {
		t.Error("a closer post should rank higher")
	}
	if Score(hour, 1000, Meetup) <= Score(48*hour, 1000, Meetup) {
		t.Error("a newer post should rank higher")
	}
	if Score(hour, 1000, HelpRequest) <= Score(hour, 1000, Meetup) {
		t.Error("a help request should rank above an otherwise equal meetup")
	}
	// A brand-new help request right next to you: 1 + 1 + urgent + help.
	if got, want := Score(0, 0, HelpRequest), 2+UrgentBoost+HelpRequestBoost; got != want {
		t.Errorf("Score(0, 0, help) = %v, want %v", got, want)
	}
}
