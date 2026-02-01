package migrations

import "database/sql"

// GetAllMigrations returns all registered migrations for the notification service
func GetAllMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Name:        "create_notifications_table",
			Description: "Creates the notifications table with indexes",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS notifications (
					id SERIAL PRIMARY KEY,
					user_id INTEGER NOT NULL,
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
			Version:     2,
			Name:        "add_websocket_connections_table",
			Description: "Creates a table to track WebSocket connections for real-time notifications",
			Up: func(db *sql.DB) error {
				query := `
				CREATE TABLE IF NOT EXISTS websocket_connections (
					id SERIAL PRIMARY KEY,
					user_id INTEGER NOT NULL,
					connection_id VARCHAR(100) UNIQUE NOT NULL,
					connected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					last_ping_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX IF NOT EXISTS idx_ws_connections_user_id ON websocket_connections (user_id);
				CREATE INDEX IF NOT EXISTS idx_ws_connections_connection_id ON websocket_connections (connection_id);
				`
				_, err := db.Exec(query)
				return err
			},
			Down: func(db *sql.DB) error {
				_, err := db.Exec("DROP TABLE IF EXISTS websocket_connections CASCADE")
				return err
			},
		},
	}
}
