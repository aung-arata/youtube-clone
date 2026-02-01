package migrations

import "database/sql"

// GetAllMigrations returns all registered migrations for the video service
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
			Name:        "create_playlists_tables",
			Description: "Creates the playlists and playlist_videos tables",
			Up: func(db *sql.DB) error {
				query := `
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
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS playlist_videos CASCADE; DROP TABLE IF EXISTS playlists CASCADE")
				return err
			},
		},
		{
			Version:     3,
			Name:        "add_video_qualities_table",
			Description: "Creates the video_qualities table for multiple quality support and transcoding",
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
					status VARCHAR(20) DEFAULT 'pending',
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
			Version:     4,
			Name:        "add_transcoding_jobs_table",
			Description: "Creates the transcoding_jobs table for tracking video transcoding",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS transcoding_jobs (
					id SERIAL PRIMARY KEY,
					video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
					target_quality VARCHAR(10) NOT NULL,
					status VARCHAR(20) DEFAULT 'pending',
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
