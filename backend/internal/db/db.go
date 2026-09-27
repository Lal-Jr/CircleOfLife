package db

import (
	"context"
	"fmt"
	"log"

	"circleoflife/pkg/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Known, publicly-documented credentials for a seeded demo account, so
// anyone reviewing this project (e.g. a recruiter) can log in immediately
// instead of having to sign up first.
const (
	DemoUserEmail    = "demo@circleoflife.app"
	DemoUserPassword = "CircleDemo123!"
)

var Pool *pgxpool.Pool

func ConnectDB(databaseURL string) {
	var err error
	Pool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err = Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	fmt.Println("Connected to PostgreSQL successfully")

	err = RunMigrations()
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func RunMigrations() error {
	ctx := context.Background()

	// 1. Ensure PostGIS is installed
	_, err := Pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS postgis;`)
	if err != nil {
		return fmt.Errorf("failed to create postgis extension: %v", err)
	}

	// 2. Users Table
	_, err = Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return fmt.Errorf("failed creating users table: %v", err)
	}

	// 3. Posts Table with Geography Point
	_, err = Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS posts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		type VARCHAR(50) NOT NULL CHECK (type IN ('help', 'meetup')),
		location geography(Point, 4326) NOT NULL,
		meetup_time TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return fmt.Errorf("failed creating posts table: %v", err)
	}

	// Spatial Index
	_, err = Pool.Exec(ctx, `
	CREATE INDEX IF NOT EXISTS idx_posts_location
	ON posts
	USING GIST(location);`)
	if err != nil {
		return fmt.Errorf("failed creating posts location index: %v", err)
	}

	// 4. Comments Table
	_, err = Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS comments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return fmt.Errorf("failed creating comments table: %v", err)
	}

	// 5. Post Likes ("Helpful" votes) - one vote per user per post
	_, err = Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS post_likes (
		post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (post_id, user_id)
	);`)
	if err != nil {
		return fmt.Errorf("failed creating post_likes table: %v", err)
	}

	fmt.Println("Database migrations applied successfully, PostGIS ready.")

	if err := seedDemoUser(ctx); err != nil {
		return fmt.Errorf("failed seeding demo user: %v", err)
	}

	return nil
}

// seedDemoUser ensures a demo account with known credentials always exists,
// with a couple of sample posts so the feed isn't empty on first login.
// Idempotent: safe to run on every startup.
func seedDemoUser(ctx context.Context) error {
	var userID string
	err := Pool.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, DemoUserEmail).Scan(&userID)
	if err == nil {
		// Already seeded.
		return nil
	}

	hashed, err := auth.HashPassword(DemoUserPassword)
	if err != nil {
		return fmt.Errorf("failed hashing demo password: %v", err)
	}

	err = Pool.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		"Demo User", DemoUserEmail, hashed,
	).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed creating demo user: %v", err)
	}

	// Sample posts near Bangalore (matches the coordinates used throughout
	// local testing/README examples) so a reviewer testing from - or
	// spoofing their location to - that area sees a populated feed.
	_, err = Pool.Exec(ctx, `
		INSERT INTO posts (user_id, title, description, type, location)
		VALUES
			($1, 'Need help carrying groceries upstairs',
			 'My elevator is out of service and I have several heavy bags. Any help would be great!',
			 'help', ST_SetSRID(ST_MakePoint(77.5946, 12.9716), 4326)),
			($1, 'Weekend park cleanup meetup',
			 'Organizing a small group to tidy up the local park this Saturday morning. All welcome!',
			 'meetup', ST_SetSRID(ST_MakePoint(77.6046, 12.9816), 4326))
	`, userID)
	if err != nil {
		return fmt.Errorf("failed seeding demo posts: %v", err)
	}

	fmt.Printf("Seeded demo account (%s) for easy review access.\n", DemoUserEmail)
	return nil
}
