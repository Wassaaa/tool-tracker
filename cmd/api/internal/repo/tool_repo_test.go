package repo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

// TestMain handles setup and teardown for the entire repo test suite
func TestMain(m *testing.M) {
	code := m.Run()

	// Cleanup shared test infrastructure
	cleanupSharedTestInfrastructure()

	os.Exit(code)
}

// TestPostgresToolRepo_CRUD tests all basic CRUD operations
func TestPostgresToolRepo_CRUD(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresToolRepo(db)

	t.Run("Create Tool", func(t *testing.T) {
		tool, err := repo.Create("Test Hammer", domain.ToolStatusInOffice)

		require.NoError(t, err)
		assert.NotNil(t, tool.ID)
		assert.Equal(t, "Test Hammer", tool.Name)
		assert.Equal(t, domain.ToolStatusInOffice, tool.Status)
		assert.Nil(t, tool.CurrentUserId)
		assert.False(t, tool.CreatedAt.IsZero())
		assert.False(t, tool.UpdatedAt.IsZero())
	})

	t.Run("Get Tool", func(t *testing.T) {
		// Create first
		created, err := repo.Create("Test Drill", domain.ToolStatusInOffice)
		require.NoError(t, err)

		// Get it back
		retrieved, err := repo.Get(*created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Name, retrieved.Name)
		assert.Equal(t, created.Status, retrieved.Status)
	})

	t.Run("Update Tool", func(t *testing.T) {
		// Create first
		created, err := repo.Create("Original Name", domain.ToolStatusInOffice)
		require.NoError(t, err)

		originalUpdatedAt := created.UpdatedAt

		// Update it
		created.Name = "Updated Name"
		created.Status = domain.ToolStatusMaintenance

		updated, err := repo.Update(created)
		require.NoError(t, err)

		assert.Equal(t, "Updated Name", updated.Name)
		assert.Equal(t, domain.ToolStatusMaintenance, updated.Status)
		// Database trigger should update the timestamp
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("Delete Tool", func(t *testing.T) {
		// Create first
		created, err := repo.Create("To Delete", domain.ToolStatusInOffice)
		require.NoError(t, err)

		// Delete it
		err = repo.Delete(*created.ID)
		require.NoError(t, err)

		// Should not be found
		_, err = repo.Get(*created.ID)
		assert.Error(t, err)
	})

	t.Run("Count Tools", func(t *testing.T) {
		// Get initial count
		initialCount, err := repo.Count()
		require.NoError(t, err)

		// Add tools
		_, err = repo.Create("Tool 1", domain.ToolStatusInOffice)
		require.NoError(t, err)

		_, err = repo.Create("Tool 2", domain.ToolStatusMaintenance)
		require.NoError(t, err)

		// Count should increase
		newCount, err := repo.Count()
		require.NoError(t, err)
		assert.Equal(t, initialCount+2, newCount)
	})
}

// TestPostgresToolRepo_PostgreSQLFeatures tests database-specific features
func TestPostgresToolRepo_PostgreSQLFeatures(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresToolRepo(db)

	t.Run("UUID Primary Keys", func(t *testing.T) {
		tool1, err := repo.Create("Tool 1", domain.ToolStatusInOffice)
		require.NoError(t, err)

		tool2, err := repo.Create("Tool 2", domain.ToolStatusInOffice)
		require.NoError(t, err)

		// UUIDs should be different
		assert.NotEqual(t, tool1.ID, tool2.ID)

		// UUIDs should be valid format (36 chars with dashes)
		assert.Len(t, *tool1.ID, 36)
		assert.Contains(t, *tool1.ID, "-")
	})

	t.Run("Database Triggers Work", func(t *testing.T) {
		// Create user first (needed for checkout)
		userID := createTestUser(t, db, "Test User", "test@example.com", domain.UserRoleEmployee)

		// Create tool
		tool, err := repo.Create("Trigger Test", domain.ToolStatusInOffice)
		require.NoError(t, err)
		assert.Nil(t, tool.LastCheckedOutAt)

		// Checkout tool (assign to user)
		tool.CurrentUserId = &userID
		tool.Status = domain.ToolStatusCheckedOut

		updated, err := repo.Update(tool)
		require.NoError(t, err)

		// Database trigger should set LastCheckedOutAt
		assert.NotNil(t, updated.LastCheckedOutAt)
		assert.True(t, updated.LastCheckedOutAt.After(tool.CreatedAt))
	})

	t.Run("List by Status", func(t *testing.T) {
		// Create tools with different statuses
		_, err := repo.Create("Available 1", domain.ToolStatusInOffice)
		require.NoError(t, err)

		_, err = repo.Create("Available 2", domain.ToolStatusInOffice)
		require.NoError(t, err)

		_, err = repo.Create("Maintenance Tool", domain.ToolStatusMaintenance)
		require.NoError(t, err)

		// Query by status
		available, err := repo.ListByStatus(domain.ToolStatusInOffice, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(available), 2)

		maintenance, err := repo.ListByStatus(domain.ToolStatusMaintenance, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(maintenance), 1)
	})
}

// TestPostgresToolRepo_ErrorCases tests error handling
func TestPostgresToolRepo_ErrorCases(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresToolRepo(db)

	t.Run("Get Non-existent Tool", func(t *testing.T) {
		_, err := repo.Get("00000000-0000-0000-0000-000000000000")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool not found")
	})

	t.Run("Get Invalid UUID", func(t *testing.T) {
		_, err := repo.Get("invalid-uuid")
		assert.Error(t, err)
	})

	t.Run("Delete Non-existent Tool", func(t *testing.T) {
		// This should error because the tool doesn't exist
		err := repo.Delete("00000000-0000-0000-0000-000000000000")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool not found")
	})
}
