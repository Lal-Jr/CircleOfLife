package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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

	fmt.Println("Database migrations applied successfully, PostGIS ready.")
	return nil
}
