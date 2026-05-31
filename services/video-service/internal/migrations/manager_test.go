package migrations

import (
	"testing"
)

func TestNewMigrationManager(t *testing.T) {
	m := NewMigrationManager(nil)
	if m == nil {
		t.Fatal("expected non-nil MigrationManager")
	}
	if m.migrations == nil {
		t.Error("expected migrations slice to be initialized")
	}
	if len(m.migrations) != 0 {
		t.Errorf("expected 0 migrations, got %d", len(m.migrations))
	}
}

func TestRegister_AddsMigration(t *testing.T) {
	m := NewMigrationManager(nil)
	m.Register(Migration{
		Version:     1,
		Name:        "create_videos_table",
		Description: "Creates the videos table",
	})

	if len(m.migrations) != 1 {
		t.Errorf("expected 1 migration, got %d", len(m.migrations))
	}
	if m.migrations[0].Version != 1 {
		t.Errorf("expected version 1, got %d", m.migrations[0].Version)
	}
	if m.migrations[0].Name != "create_videos_table" {
		t.Errorf("unexpected migration name: %s", m.migrations[0].Name)
	}
}

func TestRegister_MultipleInOrder(t *testing.T) {
	m := NewMigrationManager(nil)
	m.Register(Migration{Version: 3, Name: "add_index"})
	m.Register(Migration{Version: 1, Name: "create_table"})
	m.Register(Migration{Version: 2, Name: "add_column"})

	if len(m.migrations) != 3 {
		t.Fatalf("expected 3 migrations, got %d", len(m.migrations))
	}
	versions := map[int]bool{}
	for _, mig := range m.migrations {
		versions[mig.Version] = true
	}
	for _, v := range []int{1, 2, 3} {
		if !versions[v] {
			t.Errorf("migration version %d not found after Register", v)
		}
	}
}

func TestMigration_StructFields(t *testing.T) {
	m := Migration{
		Version:     10,
		Name:        "add_transcoding_jobs",
		Description: "Adds transcoding jobs table",
	}

	if m.Version != 10 {
		t.Errorf("expected Version 10, got %d", m.Version)
	}
	if m.Name != "add_transcoding_jobs" {
		t.Errorf("expected Name add_transcoding_jobs, got %q", m.Name)
	}
	if m.Up != nil {
		t.Error("expected Up to be nil when not set")
	}
	if m.Down != nil {
		t.Error("expected Down to be nil when not set")
	}
}
