package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	// Get database credentials from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "user_service_db")

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
	`

	_, err := db.Exec(query)
	return err
}
