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
	dbname := getEnv("DB_NAME", "history_service_db")

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
	CREATE TABLE IF NOT EXISTS watch_history (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		video_id INTEGER NOT NULL,
		watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, video_id)
	);

	CREATE INDEX IF NOT EXISTS idx_watch_history_user_id ON watch_history (user_id);
	CREATE INDEX IF NOT EXISTS idx_watch_history_watched_at ON watch_history (watched_at DESC);
	`

	_, err := db.Exec(query)
	return err
}
