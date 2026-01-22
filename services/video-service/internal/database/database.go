package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	// Get database credentials from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "video_service_db")

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)              // Maximum number of open connections
	db.SetMaxIdleConns(5)               // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection (5 minutes)

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func runMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS videos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		url VARCHAR(500) NOT NULL,
		thumbnail VARCHAR(500),
		channel_name VARCHAR(100) NOT NULL,
		channel_avatar VARCHAR(500),
		views INTEGER DEFAULT 0,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		category VARCHAR(50) DEFAULT 'General',
		duration VARCHAR(20),
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes for search performance
	CREATE INDEX IF NOT EXISTS idx_videos_title ON videos USING gin(to_tsvector('english', title));
	CREATE INDEX IF NOT EXISTS idx_videos_description ON videos USING gin(to_tsvector('english', description));
	CREATE INDEX IF NOT EXISTS idx_videos_channel_name ON videos (channel_name);
	CREATE INDEX IF NOT EXISTS idx_videos_uploaded_at ON videos (uploaded_at DESC);
	CREATE INDEX IF NOT EXISTS idx_videos_category ON videos (category);

	CREATE TABLE IF NOT EXISTS playlists (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS playlist_videos (
		id SERIAL PRIMARY KEY,
		playlist_id INTEGER REFERENCES playlists(id) ON DELETE CASCADE,
		video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
		position INTEGER NOT NULL DEFAULT 0,
		added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(playlist_id, video_id)
	);

	CREATE INDEX IF NOT EXISTS idx_playlists_user_id ON playlists (user_id);
	CREATE INDEX IF NOT EXISTS idx_playlist_videos_playlist_id ON playlist_videos (playlist_id);
	CREATE INDEX IF NOT EXISTS idx_playlist_videos_position ON playlist_videos (position);
	`

	_, err := db.Exec(query)
	return err
}
