// Package infrastructure implements the Community context's ports with Postgres/PostGIS,
// Redis and bluemonday.
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"circleoflife/internal/community/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Posts stores the Post aggregate in Postgres, with locations as PostGIS geography points.
type Posts struct {
	db *pgxpool.Pool
}

func NewPosts(db *pgxpool.Pool) *Posts { return &Posts{db: db} }

var _ domain.Posts = (*Posts)(nil)

func (r *Posts) Get(ctx context.Context, id domain.PostID) (*domain.Post, error) {
	var (
		author, title, description, kind string
		lat, lng                         float64
		meetup                           *time.Time
		created                          time.Time
	)
	err := r.db.QueryRow(ctx, `
		SELECT user_id, title, description, type, ST_Y(location::geometry), ST_X(location::geometry), meetup_time, created_at
		FROM posts WHERE id = $1`, string(id)).Scan(&author, &title, &description, &kind, &lat, &lng, &meetup, &created)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPostNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("loading post: %w", err)
	}
	at, err := domain.NewLocation(lat, lng)
	if err != nil {
		return nil, err
	}
	return domain.RehydratePost(id, domain.MemberID(author), title, description, domain.Kind(kind), at, meetup, created), nil
}

func (r *Posts) Add(ctx context.Context, p *domain.Post) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO posts (id, user_id, title, description, type, location, meetup_time, created_at)
		VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326), $8, $9)`,
		string(p.ID()), string(p.Author()), p.Title(), p.Description(), string(p.Kind()),
		p.Location().Lng(), p.Location().Lat(), p.MeetupTime(), p.CreatedAt())
	if err != nil {
		return fmt.Errorf("saving post: %w", err)
	}
	return nil
}

func (r *Posts) AddComment(ctx context.Context, c domain.Comment) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		string(c.ID()), string(c.Post()), string(c.Author()), c.Content(), c.CreatedAt())
	if err != nil {
		return fmt.Errorf("saving comment: %w", err)
	}
	return nil
}

func (r *Posts) ToggleHelpfulVote(ctx context.Context, post domain.PostID, voter domain.MemberID) (bool, int, error) {
	var voted bool
	err := pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		// Delete first: if a vote was there, this takes it back; otherwise add one.
		tag, err := tx.Exec(ctx, `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`, string(post), string(voter))
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			if _, err := tx.Exec(ctx, `INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, string(post), string(voter)); err != nil {
				return err
			}
			voted = true
		}
		return nil
	})
	if err != nil {
		return false, 0, fmt.Errorf("toggling helpful vote: %w", err)
	}
	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM post_likes WHERE post_id = $1`, string(post)).Scan(&total); err != nil {
		return false, 0, fmt.Errorf("counting helpful votes: %w", err)
	}
	return voted, total, nil
}
