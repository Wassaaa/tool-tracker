package database

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
)

// Embed all migration files into the binary at compile time
// This means the SQL files become part of your executable
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(db *sql.DB) error {
	// Step 1: Create the tracking table
	err := createMigrationsTable(db)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Step 2: Get all migration files
	migrationNames, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Step 3: Apply each migration that hasn't been applied yet
	for _, migrationName := range migrationNames {
		err := applyMigration(db, migrationName)
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migrationName, err)
		}
	}

	return nil
}

// Creates the table that tracks which migrations have been applied
func createMigrationsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )
    `
	_, err := db.Exec(query)
	return err
}

// Reads all .sql files from the migrations directory
func getMigrationFiles() ([]string, error) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	var migrationNames []string
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".sql" {
			migrationNames = append(migrationNames, entry.Name())
		}
	}
	sort.Strings(migrationNames)

	return migrationNames, nil
}

// Applies a single migration if it hasn't been applied yet
func applyMigration(db *sql.DB, migrationName string) error {
	alreadyApplied, err := isMigrationApplied(db, migrationName)
	if err != nil {
		return err
	}

	if alreadyApplied {
		fmt.Printf("Migration %s already applied, skipping\n", migrationName)
		return nil
	}

	// Read the SQL file content
	content, err := migrationFiles.ReadFile("migrations/" + migrationName)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute the SQL
	fmt.Printf("Applying migration: %s\n", migrationName)
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	// Mark this migration as applied
	err = markMigrationAsApplied(db, migrationName)
	if err != nil {
		return fmt.Errorf("failed to mark migration as applied: %w", err)
	}

	fmt.Printf("Successfully applied migration: %s\n", migrationName)
	return nil
}

// Checks if a migration has already been applied
func isMigrationApplied(db *sql.DB, migrationName string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM schema_migrations WHERE version = $1"
	err := db.QueryRow(query, migrationName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Records that a migration has been applied
func markMigrationAsApplied(db *sql.DB, migrationName string) error {
	query := "INSERT INTO schema_migrations (version) VALUES ($1)"
	_, err := db.Exec(query, migrationName)
	return err
}

// Utility function to get migration status (useful for debugging)
func GetMigrationStatus(db *sql.DB) ([]MigrationStatus, error) {
	migrationNames, err := getMigrationFiles()
	if err != nil {
		return nil, err
	}

	var statuses []MigrationStatus
	for _, name := range migrationNames {
		applied, err := isMigrationApplied(db, name)
		if err != nil {
			return nil, err
		}

		status := MigrationStatus{
			Name:    name,
			Applied: applied,
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

type MigrationStatus struct {
	Name    string
	Applied bool
}
