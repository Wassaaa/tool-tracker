package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

// TestPostgresUserRepo_CRUD tests all basic CRUD operations
func TestPostgresUserRepo_CRUD(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresUserRepo(db)

	t.Run("Create User", func(t *testing.T) {
		user, err := repo.Create("John Doe", "john@example.com", domain.UserRoleEmployee)

		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, domain.UserRoleEmployee, user.Role)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run("Get User", func(t *testing.T) {
		// Create first
		created, err := repo.Create("Jane Smith", "jane@example.com", domain.UserRoleManager)
		require.NoError(t, err)

		// Get it back
		retrieved, err := repo.Get(created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Name, retrieved.Name)
		assert.Equal(t, created.Email, retrieved.Email)
		assert.Equal(t, created.Role, retrieved.Role)
	})

	t.Run("Get User by Email", func(t *testing.T) {
		// Create first
		created, err := repo.Create("Bob Wilson", "bob@example.com", domain.UserRoleAdmin)
		require.NoError(t, err)

		// Get by email
		retrieved, err := repo.GetByEmail("bob@example.com")
		require.NoError(t, err)

		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Name, retrieved.Name)
		assert.Equal(t, created.Email, retrieved.Email)
	})

	t.Run("Update User", func(t *testing.T) {
		// Create first
		created, err := repo.Create("Alice Brown", "alice@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		originalUpdatedAt := created.UpdatedAt

		// Update it
		updated, err := repo.Update(created.ID, "Alice Johnson", "alice.johnson@example.com", domain.UserRoleManager)
		require.NoError(t, err)

		assert.Equal(t, "Alice Johnson", updated.Name)
		assert.Equal(t, "alice.johnson@example.com", updated.Email)
		assert.Equal(t, domain.UserRoleManager, updated.Role)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("Delete User", func(t *testing.T) {
		// Create first
		created, err := repo.Create("Charlie Davis", "charlie@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		// Delete it
		err = repo.Delete(created.ID)
		require.NoError(t, err)

		// Should not be found
		_, err = repo.Get(created.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Count Users", func(t *testing.T) {
		// Get initial count
		initialCount, err := repo.Count()
		require.NoError(t, err)

		// Add users
		_, err = repo.Create("User 1", "user1@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		_, err = repo.Create("User 2", "user2@example.com", domain.UserRoleAdmin)
		require.NoError(t, err)

		// Count should increase
		newCount, err := repo.Count()
		require.NoError(t, err)
		assert.Equal(t, initialCount+2, newCount)
	})
}

// TestPostgresUserRepo_QueryFeatures tests specialized query operations
func TestPostgresUserRepo_QueryFeatures(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresUserRepo(db)

	t.Run("List Users", func(t *testing.T) {
		// Create test users
		_, err := repo.Create("User A", "usera@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		_, err = repo.Create("User B", "userb@example.com", domain.UserRoleManager)
		require.NoError(t, err)

		_, err = repo.Create("User C", "userc@example.com", domain.UserRoleAdmin)
		require.NoError(t, err)

		// List with pagination
		users, err := repo.List(2, 0)
		require.NoError(t, err)
		assert.Len(t, users, 2)

		// List next page
		moreUsers, err := repo.List(2, 2)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(moreUsers), 1)

		// Should be ordered by created_at DESC (newest first)
		if len(users) >= 2 {
			assert.True(t, users[0].CreatedAt.After(users[1].CreatedAt) || users[0].CreatedAt.Equal(users[1].CreatedAt))
		}
	})

	t.Run("List by Role", func(t *testing.T) {
		// Create users with different roles
		_, err := repo.Create("Employee 1", "emp1@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		_, err = repo.Create("Employee 2", "emp2@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		_, err = repo.Create("Manager 1", "mgr1@example.com", domain.UserRoleManager)
		require.NoError(t, err)

		_, err = repo.Create("Admin 1", "admin1@example.com", domain.UserRoleAdmin)
		require.NoError(t, err)

		// Query by role
		employees, err := repo.ListByRole(domain.UserRoleEmployee, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(employees), 2)
		for _, user := range employees {
			assert.Equal(t, domain.UserRoleEmployee, user.Role)
		}

		managers, err := repo.ListByRole(domain.UserRoleManager, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(managers), 1)
		for _, user := range managers {
			assert.Equal(t, domain.UserRoleManager, user.Role)
		}

		admins, err := repo.ListByRole(domain.UserRoleAdmin, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(admins), 1)
		for _, user := range admins {
			assert.Equal(t, domain.UserRoleAdmin, user.Role)
		}
	})

	t.Run("UUID Primary Keys", func(t *testing.T) {
		user1, err := repo.Create("User 1", "uuid1@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		user2, err := repo.Create("User 2", "uuid2@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		// UUIDs should be different
		assert.NotEqual(t, user1.ID, user2.ID)

		// UUIDs should be valid format (36 chars with dashes)
		assert.Len(t, user1.ID, 36)
		assert.Contains(t, user1.ID, "-")
	})
}

// TestPostgresUserRepo_ErrorCases tests error handling
func TestPostgresUserRepo_ErrorCases(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresUserRepo(db)

	t.Run("Get Non-existent User", func(t *testing.T) {
		_, err := repo.Get("00000000-0000-0000-0000-000000000000")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Get Invalid UUID", func(t *testing.T) {
		_, err := repo.Get("invalid-uuid")
		assert.Error(t, err)
	})

	t.Run("Get by Non-existent Email", func(t *testing.T) {
		_, err := repo.GetByEmail("nonexistent@example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Update Non-existent User", func(t *testing.T) {
		_, err := repo.Update("00000000-0000-0000-0000-000000000000", "Updated Name", "updated@example.com", domain.UserRoleEmployee)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Delete Non-existent User", func(t *testing.T) {
		err := repo.Delete("00000000-0000-0000-0000-000000000000")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Create User with Duplicate Email", func(t *testing.T) {
		// Create first user
		_, err := repo.Create("First User", "duplicate@example.com", domain.UserRoleEmployee)
		require.NoError(t, err)

		// Try to create another user with same email (should fail due to unique constraint)
		_, err = repo.Create("Second User", "duplicate@example.com", domain.UserRoleEmployee)
		assert.Error(t, err)
		// The exact error message will depend on the database constraint
	})
}
