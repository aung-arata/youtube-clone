package migrations

import "database/sql"

// GetAllMigrations returns all registered migrations for the backend
func GetAllMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Name:        "create_videos_table",
			Description: "Creates the videos table with indexes",
			Up: func(db *sql.DB) error {
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

				CREATE INDEX IF NOT EXISTS idx_videos_title ON videos USING gin(to_tsvector('english', title));
				CREATE INDEX IF NOT EXISTS idx_videos_description ON videos USING gin(to_tsvector('english', description));
				CREATE INDEX IF NOT EXISTS idx_videos_channel_name ON videos (channel_name);
				CREATE INDEX IF NOT EXISTS idx_videos_uploaded_at ON videos (uploaded_at DESC);
				CREATE INDEX IF NOT EXISTS idx_videos_category ON videos (category);
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS videos CASCADE")
				return err
			},
		},
		{
			Version:     2,
			Name:        "create_users_table",
			Description: "Creates the users table",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS users (
					id SERIAL PRIMARY KEY,
					username VARCHAR(50) UNIQUE NOT NULL,
					email VARCHAR(100) UNIQUE NOT NULL,
					password VARCHAR(255) NOT NULL DEFAULT '',
					avatar VARCHAR(500),
					role VARCHAR(50) NOT NULL DEFAULT 'user',
					plan_id INTEGER,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS users CASCADE")
				return err
			},
		},
		{
			Version:     3,
			Name:        "create_plans_table",
			Description: "Creates the subscription plans table",
			Up: func(db *sql.DB) error {
				query := `
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
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS plans CASCADE")
				return err
			},
		},
		{
			Version:     4,
			Name:        "create_comments_table",
			Description: "Creates the comments table",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS comments (
					id SERIAL PRIMARY KEY,
					video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
					user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
					content TEXT NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS comments CASCADE")
				return err
			},
		},
		{
			Version:     5,
			Name:        "create_watch_history_table",
			Description: "Creates the watch history table",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS watch_history (
					id SERIAL PRIMARY KEY,
					user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
					video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
					watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(user_id, video_id)
				);

				CREATE INDEX IF NOT EXISTS idx_watch_history_user_id ON watch_history (user_id);
				CREATE INDEX IF NOT EXISTS idx_watch_history_watched_at ON watch_history (watched_at DESC);
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS watch_history CASCADE")
				return err
			},
		},
		{
			Version:     6,
			Name:        "create_subscriptions_table",
			Description: "Creates the channel subscriptions table",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS subscriptions (
					id SERIAL PRIMARY KEY,
					user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
					channel_name VARCHAR(100) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(user_id, channel_name)
				);

				CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
				CREATE INDEX IF NOT EXISTS idx_subscriptions_channel_name ON subscriptions (channel_name);
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS subscriptions CASCADE")
				return err
			},
		},
		{
			Version:     7,
			Name:        "create_playlists_tables",
			Description: "Creates the playlists and playlist_videos tables",
			Up: func(db *sql.DB) error {
				query := `
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
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS playlist_videos CASCADE; DROP TABLE IF EXISTS playlists CASCADE")
				return err
			},
		},
		{
			Version:     8,
			Name:        "create_notifications_table",
			Description: "Creates the notifications table",
			Up: func(db *sql.DB) error {
				query := `
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
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS notifications CASCADE")
				return err
			},
		},
		{
			Version:     9,
			Name:        "add_video_qualities_table",
			Description: "Creates the video_qualities table for multiple quality support",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS video_qualities (
					id SERIAL PRIMARY KEY,
					video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
					quality VARCHAR(10) NOT NULL,
					url VARCHAR(500) NOT NULL,
					bitrate INTEGER,
					width INTEGER,
					height INTEGER,
					format VARCHAR(20) DEFAULT 'mp4',
					file_size BIGINT,
					status VARCHAR(20) NOT NULL DEFAULT 'pending',
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(video_id, quality)
				);

				CREATE INDEX IF NOT EXISTS idx_video_qualities_video_id ON video_qualities (video_id);
				CREATE INDEX IF NOT EXISTS idx_video_qualities_status ON video_qualities (status);
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS video_qualities CASCADE")
				return err
			},
		},
		{
			Version:     10,
			Name:        "add_transcoding_jobs_table",
			Description: "Creates the transcoding_jobs table for tracking video transcoding",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS transcoding_jobs (
					id SERIAL PRIMARY KEY,
					video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
					target_quality VARCHAR(10) NOT NULL,
					status VARCHAR(20) NOT NULL DEFAULT 'pending',
					progress INTEGER DEFAULT 0,
					error_message TEXT,
					started_at TIMESTAMP,
					completed_at TIMESTAMP,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(video_id, target_quality)
				);

				CREATE INDEX IF NOT EXISTS idx_transcoding_jobs_video_id ON transcoding_jobs (video_id);
				CREATE INDEX IF NOT EXISTS idx_transcoding_jobs_status ON transcoding_jobs (status);
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS transcoding_jobs CASCADE")
				return err
			},
		},
	}
}
