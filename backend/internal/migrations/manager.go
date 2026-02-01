package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Name        string
	Description string
	Up          func(db *sql.DB) error
	Down        func(db *sql.DB) error
}

// MigrationManager manages database migrations
type MigrationManager struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// Register adds a migration to the manager
func (m *MigrationManager) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// EnsureMigrationTable creates the migrations tracking table if it doesn't exist
func (m *MigrationManager) EnsureMigrationTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		execution_time_ms INTEGER
	);
	`
	_, err := m.db.Exec(query)
	return err
}

// GetAppliedMigrations returns a list of applied migration versions
func (m *MigrationManager) GetAppliedMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)

	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return applied, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return applied, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// GetCurrentVersion returns the highest applied migration version
func (m *MigrationManager) GetCurrentVersion() (int, error) {
	var version int
	err := m.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	return version, err
}

// MigrateUp runs all pending migrations
func (m *MigrationManager) MigrateUp() error {
	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	for _, migration := range m.migrations {
		if applied[migration.Version] {
			continue
		}

		log.Printf("Applying migration %d: %s", migration.Version, migration.Name)

		start := time.Now()
		if err := migration.Up(m.db); err != nil {
			return fmt.Errorf("failed to apply migration %d (%s): %w", migration.Version, migration.Name, err)
		}
		executionTime := time.Since(start).Milliseconds()

		// Record the migration
		_, err := m.db.Exec(
			"INSERT INTO schema_migrations (version, name, description, execution_time_ms) VALUES ($1, $2, $3, $4)",
			migration.Version, migration.Name, migration.Description, executionTime,
		)
		if err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		log.Printf("Applied migration %d in %dms", migration.Version, executionTime)
	}

	return nil
}

// MigrateDown rolls back the last migration
func (m *MigrationManager) MigrateDown() error {
	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	currentVersion, err := m.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if currentVersion == 0 {
		log.Println("No migrations to roll back")
		return nil
	}

	// Find the migration to roll back
	var migrationToRollback *Migration
	for i := range m.migrations {
		if m.migrations[i].Version == currentVersion {
			migrationToRollback = &m.migrations[i]
			break
		}
	}

	if migrationToRollback == nil {
		return fmt.Errorf("migration %d not found", currentVersion)
	}

	if migrationToRollback.Down == nil {
		return fmt.Errorf("migration %d has no down function", currentVersion)
	}

	log.Printf("Rolling back migration %d: %s", currentVersion, migrationToRollback.Name)

	if err := migrationToRollback.Down(m.db); err != nil {
		return fmt.Errorf("failed to roll back migration %d: %w", currentVersion, err)
	}

	// Remove the migration record
	_, err = m.db.Exec("DELETE FROM schema_migrations WHERE version = $1", currentVersion)
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	log.Printf("Rolled back migration %d", currentVersion)
	return nil
}

// MigrateDownTo rolls back to a specific version
func (m *MigrationManager) MigrateDownTo(targetVersion int) error {
	for {
		currentVersion, err := m.GetCurrentVersion()
		if err != nil {
			return err
		}

		if currentVersion <= targetVersion {
			break
		}

		if err := m.MigrateDown(); err != nil {
			return err
		}
	}
	return nil
}

// Status prints the migration status
func (m *MigrationManager) Status() error {
	if err := m.EnsureMigrationTable(); err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	log.Println("Migration Status:")
	log.Println("================")

	for _, migration := range m.migrations {
		status := "Pending"
		if applied[migration.Version] {
			status = "Applied"
		}
		log.Printf("[%s] %d: %s", status, migration.Version, migration.Name)
	}

	return nil
}

// GetPendingMigrations returns migrations that haven't been applied
func (m *MigrationManager) GetPendingMigrations() ([]Migration, error) {
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	var pending []Migration
	for _, migration := range m.migrations {
		if !applied[migration.Version] {
			pending = append(pending, migration)
		}
	}

	// Sort by version
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Version < pending[j].Version
	})

	return pending, nil
}
