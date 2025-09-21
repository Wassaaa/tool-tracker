package repo

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/database"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"

	_ "github.com/lib/pq"
)

// Shared test infrastructure for all repo tests
var (
	sharedTestDB        *sql.DB
	sharedTestContainer *postgres.PostgresContainer
	sharedSetupOnce     sync.Once
)

// setupSharedRepoTestDB creates a single PostgreSQL container for ALL repo tests
func setupSharedRepoTestDB(t *testing.T) *sql.DB {
	sharedSetupOnce.Do(func() {
		ctx := context.Background()

		// Create a single PostgreSQL container for all repo tests
		var err error
		sharedTestContainer, err = postgres.Run(ctx, "postgres",
			postgres.WithDatabase("test_tooltracker_repo"),
			postgres.WithUsername("test_user"),
			postgres.WithPassword("test_pass"),
		)
		require.NoError(t, err)

		// Get connection string
		connStr, err := sharedTestContainer.ConnectionString(ctx, "sslmode=disable")
		require.NoError(t, err)

		// Connect to the test database
		sharedTestDB, err = sql.Open("postgres", connStr)
		require.NoError(t, err)

		// Verify connection with retry
		var pingErr error
		for i := 0; i < 10; i++ {
			pingErr = sharedTestDB.Ping()
			if pingErr == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		require.NoError(t, pingErr, "Failed to connect to shared test database")

		err = database.RunMigrations(sharedTestDB)
		require.NoError(t, err, "Failed to run migrations on shared test database")
	})

	// Clean up data before each test
	cleanupSharedTestData(t, sharedTestDB)
	return sharedTestDB
}

// cleanupSharedTestData removes all test data while preserving schema
func cleanupSharedTestData(t *testing.T, db *sql.DB) {
	// Delete in reverse order of dependencies
	tables := []string{"events", "tools", "users"}
	for _, table := range tables {
		// Skip system user (id = 1) if it exists
		query := "DELETE FROM " + table
		if table == "users" {
			query += " WHERE id != '00000000-0000-0000-0000-000000000001'"
		}
		_, err := db.Exec(query)
		require.NoError(t, err, "Failed to clean up table: "+table)
	}
}

// Helper functions for creating test data across all repo tests

// createTestUser creates a test user for testing across all repos
func createTestUser(t *testing.T, db *sql.DB, name, email string, role domain.UserRole) string {
	var userID string
	err := db.QueryRow(
		"INSERT INTO users (name, email, role) VALUES ($1, $2, $3) RETURNING id",
		name, email, role,
	).Scan(&userID)
	require.NoError(t, err)
	return userID
}

// createTestTool creates a test tool for testing across all repos
func createTestTool(t *testing.T, db *sql.DB, name string, status domain.ToolStatus) string {
	var toolID string
	err := db.QueryRow(
		"INSERT INTO tools (name, status) VALUES ($1, $2) RETURNING id",
		name, status,
	).Scan(&toolID)
	require.NoError(t, err)
	return toolID
}

// createTestEvent creates a test event for testing
func createTestEvent(t *testing.T, db *sql.DB, eventType domain.EventType, toolID, userID, actorID *string, notes string) string {
	var eventID string
	err := db.QueryRow(
		"INSERT INTO events (type, tool_id, user_id, actor_id, notes) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		eventType, toolID, userID, actorID, notes,
	).Scan(&eventID)
	require.NoError(t, err)
	return eventID
}

// createTestToolWithUser creates a tool checked out to a specific user
func createTestToolWithUser(t *testing.T, db *sql.DB, toolName string, userName, userEmail string, role domain.UserRole) (toolID, userID string) {
	userID = createTestUser(t, db, userName, userEmail, role)

	var tool_id string
	err := db.QueryRow(
		"INSERT INTO tools (name, status, current_user_id) VALUES ($1, $2, $3) RETURNING id",
		toolName, domain.ToolStatusCheckedOut, userID,
	).Scan(&tool_id)
	require.NoError(t, err)
	return tool_id, userID
}

// assertEventExists verifies an event exists with specific criteria
func assertEventExists(t *testing.T, db *sql.DB, eventType domain.EventType, toolID, userID *string) {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM events WHERE type = $1 AND ($2::uuid IS NULL OR tool_id = $2) AND ($3::uuid IS NULL OR user_id = $3 OR actor_id = $3)",
		eventType, toolID, userID,
	).Scan(&count)
	require.NoError(t, err)
	require.Greater(t, count, 0, "Expected event of type %s not found", eventType)
}

// cleanupSharedTestInfrastructure is called by TestMain to cleanup containers
func cleanupSharedTestInfrastructure() {
	ctx := context.Background()

	if sharedTestContainer != nil {
		_ = sharedTestContainer.Terminate(ctx)
	}
	if sharedTestDB != nil {
		sharedTestDB.Close()
	}
}
