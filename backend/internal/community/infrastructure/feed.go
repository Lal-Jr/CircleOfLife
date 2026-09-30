package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"circleoflife/internal/community/application"
	"circleoflife/internal/community/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Feed is the read side: projections of posts and comments joined with authors' names from
// the Identity context's users table, which shares the database.
type Feed struct {
	db *pgxpool.Pool
}

func NewFeed(db *pgxpool.Pool) *Feed { return &Feed{db: db} }

var _ application.FeedReader = (*Feed)(nil)

// distanceSQL is the viewer's distance to a post in metres; $1/$2 are the viewer's lng/lat.
const distanceSQL = `ST_Distance(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography)`

// scoreSQL is domain.Score written in SQL, so the database can rank and paginate the feed.
// It's generated from the domain's constants, so the two can't drift apart.
var scoreSQL = fmt.Sprintf(`(
	1.0 / (EXTRACT(EPOCH FROM (NOW() - p.created_at)) / 3600 + 1.0) +
	1.0 / (%[1]s / 1000 + 1.0) +
	CASE WHEN %[1]s < %[2]d THEN %[4]v WHEN %[1]s < %[3]d THEN %[5]v ELSE 0.0 END +
	CASE WHEN p.type = '%[6]s' THEN %[7]v ELSE 0.0 END
)`, distanceSQL, domain.UrgentWithinMeters, domain.NearbyWithinMeters, domain.UrgentBoost, domain.NearbyBoost, domain.HelpRequest, domain.HelpRequestBoost)

// viewColumns are the columns scanned by scanView; $viewer is the viewer's member id.
func viewColumns(viewerParam int) string {
	return fmt.Sprintf(`
		p.id, p.user_id, p.title, p.description, p.type, p.meetup_time, p.created_at,
		u.name AS author,
		ST_Y(p.location::geometry) AS lat,
		ST_X(p.location::geometry) AS lng,
		%s AS distance,
		(SELECT count(*) FROM comments c WHERE c.post_id = p.id) AS comment_count,
		(SELECT count(*) FROM post_likes pl WHERE pl.post_id = p.id) AS helpful_count,
		EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = $%d) AS liked_by_me`, distanceSQL, viewerParam)
}

func scanView(row pgx.Row) (application.PostView, error) {
	var v application.PostView
	err := row.Scan(&v.ID, &v.UserID, &v.Title, &v.Description, &v.Type, &v.MeetupTime, &v.CreatedAt,
		&v.Author, &v.Lat, &v.Lng, &v.Distance, &v.CommentCount, &v.HelpfulCount, &v.LikedByMe)
	return v, err
}

func (f *Feed) Nearby(ctx context.Context, req application.FeedRequest) ([]application.PostView, error) {
	// ST_DWithin uses the GiST index on location, so only posts inside the radius are ranked.
	query := `SELECT ` + viewColumns(6) + `
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE ST_DWithin(p.location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)
		ORDER BY ` + scoreSQL + ` DESC, distance ASC, p.created_at DESC
		LIMIT $4 OFFSET $5`

	// One extra row tells the caller whether another page exists.
	rows, err := f.db.Query(ctx, query, req.At.Lng, req.At.Lat, req.RadiusKm*1000, req.Limit+1, (req.Page-1)*req.Limit, string(req.Viewer))
	if err != nil {
		return nil, fmt.Errorf("querying nearby posts: %w", err)
	}
	defer rows.Close()

	var views []application.PostView
	for rows.Next() {
		v, err := scanView(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning post row: %w", err)
		}
		v.Priority = string(domain.ProximityOf(v.Distance))
		views = append(views, v)
	}
	return views, rows.Err()
}

func (f *Feed) Post(ctx context.Context, id domain.PostID, from application.Point, viewer domain.MemberID) (*application.PostView, error) {
	v, err := scanView(f.db.QueryRow(ctx, `SELECT `+viewColumns(4)+`
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $3`, from.Lng, from.Lat, string(id), string(viewer)))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPostNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("loading post view: %w", err)
	}
	return &v, nil
}

const commentColumns = `c.id, c.post_id, c.user_id, c.content, c.created_at, u.name AS author_name
	FROM comments c JOIN users u ON c.user_id = u.id`

func scanComment(row pgx.Row) (application.CommentView, error) {
	var c application.CommentView
	err := row.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.AuthorName)
	return c, err
}

func (f *Feed) Comments(ctx context.Context, post domain.PostID) ([]application.CommentView, error) {
	rows, err := f.db.Query(ctx, `SELECT `+commentColumns+` WHERE c.post_id = $1 ORDER BY c.created_at ASC`, string(post))
	if err != nil {
		return nil, fmt.Errorf("querying comments: %w", err)
	}
	defer rows.Close()
	var comments []application.CommentView
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning comment row: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (f *Feed) Comment(ctx context.Context, id domain.CommentID) (*application.CommentView, error) {
	c, err := scanComment(f.db.QueryRow(ctx, `SELECT `+commentColumns+` WHERE c.id = $1`, string(id)))
	if err != nil {
		return nil, fmt.Errorf("loading comment view: %w", err)
	}
	return &c, nil
}
