package infrastructure

import (
	"context"
	"math"
	"os"
	"testing"
	"time"

	"circleoflife/internal/community/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestScoreSQLMatchesDomainScore checks that the ranking the database applies is exactly
// domain.Score. It needs a PostGIS database: set TEST_DATABASE_URL to run it.
func TestScoreSQLMatchesDomainScore(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	ctx := context.Background()
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS postgis`); err != nil {
		t.Fatal(err)
	}

	viewerLat, viewerLng := 12.9716, 77.5946
	cases := []struct {
		dLat float64 // offset of the post from the viewer, in degrees of latitude
		age  time.Duration
		kind domain.Kind
	}{
		{0, 0, domain.HelpRequest},
		{0.003, 2 * time.Hour, domain.Meetup},        // ~330 m: urgent
		{0.01, 30 * time.Minute, domain.HelpRequest}, // ~1.1 km: nearby
		{0.05, 72 * time.Hour, domain.Meetup},        // ~5.5 km: distant
	}
	for _, c := range cases {
		var sqlScore, distance float64
		err := db.QueryRow(ctx, `
			WITH p AS (
				SELECT now() - make_interval(secs => $3::float8) AS created_at,
				       ST_SetSRID(ST_MakePoint($1::float8, $2::float8 + $4::float8), 4326)::geography AS location,
				       $5::text AS type
			)
			SELECT `+scoreSQL+`, `+distanceSQL+` FROM p`,
			viewerLng, viewerLat, c.age.Seconds(), c.dLat, string(c.kind)).Scan(&sqlScore, &distance)
		if err != nil {
			t.Fatal(err)
		}
		want := domain.Score(c.age, distance, c.kind)
		if math.Abs(sqlScore-want) > 1e-6 {
			t.Errorf("%+v: SQL score %v, domain score %v (distance %.0f m)", c, sqlScore, want, distance)
		}
	}
}
