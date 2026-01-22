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
	dbname := getEnv("DB_NAME", "youtube_clone")

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

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		avatar VARCHAR(500),
		plan_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS plans (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price DECIMAL(10, 2) NOT NULL DEFAULT 0,
		max_video_quality VARCHAR(10) DEFAULT '1080p',
		max_uploads_per_month INTEGER DEFAULT 0,
		ads_free BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Add foreign key constraint for users.plan_id
	DO $$ 
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_constraint WHERE conname = 'users_plan_id_fkey'
		) THEN
			ALTER TABLE users ADD CONSTRAINT users_plan_id_fkey 
			FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE SET NULL;
		END IF;
	END $$;

	-- Insert default plans if they don't exist
	INSERT INTO plans (id, name, price, max_video_quality, max_uploads_per_month, ads_free)
	VALUES 
		(1, 'Free', 0, '480p', 5, FALSE),
		(2, 'Basic', 4.99, '720p', 20, FALSE),
		(3, 'Premium', 9.99, '1080p', 100, TRUE),
		(4, 'Enterprise', 19.99, '4K', -1, TRUE)
	ON CONFLICT (id) DO NOTHING;

	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS watch_history (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
		watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, video_id)
	);

	CREATE INDEX IF NOT EXISTS idx_watch_history_user_id ON watch_history (user_id);
	CREATE INDEX IF NOT EXISTS idx_watch_history_watched_at ON watch_history (watched_at DESC);

	CREATE TABLE IF NOT EXISTS subscriptions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		channel_name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, channel_name)
	);

	CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
	CREATE INDEX IF NOT EXISTS idx_subscriptions_channel_name ON subscriptions (channel_name);

	CREATE TABLE IF NOT EXISTS playlists (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
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

	CREATE TABLE IF NOT EXISTS notifications (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		type VARCHAR(50) NOT NULL,
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		link VARCHAR(500),
		is_read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications (user_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications (created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications (is_read);
	`

	_, err := db.Exec(query)
	return err
}
