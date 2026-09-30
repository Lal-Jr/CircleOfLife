package infrastructure

import (
	"context"

	"circleoflife/internal/community/application"
	"circleoflife/internal/community/domain"
	"circleoflife/internal/realtime"

	"github.com/microcosm-cc/bluemonday"
)

// PlainText strips all HTML with bluemonday's strict policy.
type PlainText struct {
	policy *bluemonday.Policy
}

func NewPlainText() PlainText { return PlainText{policy: bluemonday.StrictPolicy()} }

func (t PlainText) Clean(s string) string { return t.policy.Sanitize(s) }

var _ application.TextPolicy = PlainText{}

// LiveUpdates translates Community events into the Realtime context's messages. This is the
// only place the two contexts meet.
type LiveUpdates struct {
	bus *realtime.Bus
}

func NewLiveUpdates(bus *realtime.Bus) LiveUpdates { return LiveUpdates{bus: bus} }

var _ application.EventPublisher = LiveUpdates{}

func (l LiveUpdates) Publish(ctx context.Context, events []domain.Event) {
	for _, e := range events {
		switch e := e.(type) {
		case domain.PostPublished:
			// Detached from the request: a slow publish shouldn't be cancelled by the client
			// hanging up once it has its 201.
			l.bus.Publish(context.WithoutCancel(ctx), realtime.Event{
				Type:   realtime.PostCreated,
				PostID: string(e.Post),
				Lat:    e.At.Lat(),
				Lng:    e.At.Lng(),
			})
		}
	}
}
