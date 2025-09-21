package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

// TestPostgresEventRepo_CRUD tests all basic CRUD operations
func TestPostgresEventRepo_CRUD(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresEventRepo(db)

	// Create test data
	userID := createTestUser(t, db, "Test User", "test@example.com", domain.UserRoleEmployee)
	toolID := createTestTool(t, db, "Test Tool", domain.ToolStatusInOffice)
	actorID := createTestUser(t, db, "Test Actor", "actor@example.com", domain.UserRoleManager)

	t.Run("Create Event - Tool Created", func(t *testing.T) {
		event, err := repo.Create(
			domain.EventTypeToolCreated,
			&toolID,
			nil,
			&actorID,
			"Tool was created",
			nil,
		)

		require.NoError(t, err)
		assert.NotEmpty(t, event.ID)
		assert.Equal(t, domain.EventTypeToolCreated, event.Type)
		assert.Equal(t, &toolID, event.ToolID)
		assert.Nil(t, event.UserID)
		assert.Equal(t, &actorID, event.ActorID)
		assert.Equal(t, "Tool was created", event.Notes)
		assert.Nil(t, event.Metadata)
		assert.False(t, event.CreatedAt.IsZero())
	})

	t.Run("Create Event - Tool Checked Out", func(t *testing.T) {
		metadata := `{"location": "workshop", "expected_return": "2024-01-15"}`
		event, err := repo.Create(
			domain.EventTypeToolCheckedOut,
			&toolID,
			&userID,
			&actorID,
			"Tool checked out to user",
			&metadata,
		)

		require.NoError(t, err)
		assert.Equal(t, domain.EventTypeToolCheckedOut, event.Type)
		assert.Equal(t, &toolID, event.ToolID)
		assert.Equal(t, &userID, event.UserID)
		assert.Equal(t, &actorID, event.ActorID)
		assert.Equal(t, &metadata, event.Metadata)
	})

	t.Run("Get Event", func(t *testing.T) {
		// Create first
		created, err := repo.Create(
			domain.EventTypeUserCreated,
			nil,
			&userID,
			&actorID,
			"User was created",
			nil,
		)
		require.NoError(t, err)

		// Get it back
		retrieved, err := repo.Get(created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Type, retrieved.Type)
		assert.Equal(t, created.UserID, retrieved.UserID)
		assert.Equal(t, created.ActorID, retrieved.ActorID)
		assert.Equal(t, created.Notes, retrieved.Notes)
	})

	t.Run("Count Events", func(t *testing.T) {
		// Get initial count
		initialCount, err := repo.Count()
		require.NoError(t, err)

		// Add events
		_, err = repo.Create(domain.EventTypeToolMaintenance, &toolID, nil, &actorID, "Maintenance started", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeUserUpdated, nil, &userID, &actorID, "User updated", nil)
		require.NoError(t, err)

		// Count should increase
		newCount, err := repo.Count()
		require.NoError(t, err)
		assert.Equal(t, initialCount+2, newCount)
	})
}

// TestPostgresEventRepo_QueryFeatures tests specialized query operations
func TestPostgresEventRepo_QueryFeatures(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresEventRepo(db)

	// Create test data
	user1ID := createTestUser(t, db, "User 1", "user1@example.com", domain.UserRoleEmployee)
	user2ID := createTestUser(t, db, "User 2", "user2@example.com", domain.UserRoleEmployee)
	tool1ID := createTestTool(t, db, "Tool 1", domain.ToolStatusInOffice)
	tool2ID := createTestTool(t, db, "Tool 2", domain.ToolStatusInOffice)
	actorID := createTestUser(t, db, "Actor", "actor@example.com", domain.UserRoleAdmin)

	t.Run("List Events", func(t *testing.T) {
		// Create test events
		_, err := repo.Create(domain.EventTypeToolCreated, &tool1ID, nil, &actorID, "Tool 1 created", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolCreated, &tool2ID, nil, &actorID, "Tool 2 created", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeUserCreated, nil, &user1ID, &actorID, "User 1 created", nil)
		require.NoError(t, err)

		// List with pagination
		events, err := repo.List(2, 0)
		require.NoError(t, err)
		assert.Len(t, events, 2)

		// Should be ordered by created_at DESC (newest first)
		if len(events) >= 2 {
			assert.True(t, events[0].CreatedAt.After(events[1].CreatedAt) || events[0].CreatedAt.Equal(events[1].CreatedAt))
		}
	})

	t.Run("List by Type", func(t *testing.T) {
		// Create events of different types
		_, err := repo.Create(domain.EventTypeToolCheckedOut, &tool1ID, &user1ID, &actorID, "Checkout 1", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolCheckedOut, &tool2ID, &user2ID, &actorID, "Checkout 2", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolCheckedIn, &tool1ID, &user1ID, &actorID, "Checkin 1", nil)
		require.NoError(t, err)

		// Query by type
		checkouts, err := repo.ListByType(domain.EventTypeToolCheckedOut, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(checkouts), 2)
		for _, event := range checkouts {
			assert.Equal(t, domain.EventTypeToolCheckedOut, event.Type)
		}

		checkins, err := repo.ListByType(domain.EventTypeToolCheckedIn, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(checkins), 1)
		for _, event := range checkins {
			assert.Equal(t, domain.EventTypeToolCheckedIn, event.Type)
		}
	})

	t.Run("List by Tool", func(t *testing.T) {
		// Create events for specific tool
		_, err := repo.Create(domain.EventTypeToolCreated, &tool1ID, nil, &actorID, "Tool created", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolCheckedOut, &tool1ID, &user1ID, &actorID, "Tool checked out", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolCheckedIn, &tool1ID, &user1ID, &actorID, "Tool checked in", nil)
		require.NoError(t, err)

		// Create event for different tool
		_, err = repo.Create(domain.EventTypeToolCreated, &tool2ID, nil, &actorID, "Other tool created", nil)
		require.NoError(t, err)

		// Query events for tool1
		tool1Events, err := repo.ListByTool(tool1ID, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tool1Events), 3)
		for _, event := range tool1Events {
			assert.Equal(t, &tool1ID, event.ToolID)
		}
	})

	t.Run("List by User", func(t *testing.T) {
		// Create events where user1 is involved (as user or actor)
		_, err := repo.Create(domain.EventTypeUserCreated, nil, &user1ID, &actorID, "User created", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolCheckedOut, &tool1ID, &user1ID, &actorID, "User checked out tool", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolUpdated, &tool1ID, nil, &user1ID, "User updated tool", nil)
		require.NoError(t, err)

		// Create event for different user
		_, err = repo.Create(domain.EventTypeUserCreated, nil, &user2ID, &actorID, "Other user created", nil)
		require.NoError(t, err)

		// Query events for user1 (should include events where they are user_id or actor_id)
		user1Events, err := repo.ListByUser(user1ID, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(user1Events), 3)
		for _, event := range user1Events {
			// User should be either the subject (user_id) or the actor (actor_id)
			assert.True(t,
				(event.UserID != nil && *event.UserID == user1ID) ||
					(event.ActorID != nil && *event.ActorID == user1ID),
			)
		}
	})

	t.Run("List with Filter", func(t *testing.T) {
		// Create various events
		_, err := repo.Create(domain.EventTypeToolCheckedOut, &tool1ID, &user1ID, &actorID, "Filter test 1", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolCheckedIn, &tool1ID, &user1ID, &actorID, "Filter test 2", nil)
		require.NoError(t, err)

		_, err = repo.Create(domain.EventTypeToolCheckedOut, &tool2ID, &user2ID, &actorID, "Filter test 3", nil)
		require.NoError(t, err)

		// Filter by type only
		eventType := domain.EventTypeToolCheckedOut
		typeFilter := EventFilter{Type: &eventType}
		checkoutEvents, err := repo.ListWithFilter(typeFilter, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(checkoutEvents), 2)
		for _, event := range checkoutEvents {
			assert.Equal(t, domain.EventTypeToolCheckedOut, event.Type)
		}

		// Filter by tool only
		toolFilter := EventFilter{ToolID: &tool1ID}
		tool1Events, err := repo.ListWithFilter(toolFilter, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tool1Events), 2)
		for _, event := range tool1Events {
			assert.Equal(t, &tool1ID, event.ToolID)
		}

		// Filter by user only
		userFilter := EventFilter{UserID: &user1ID}
		user1Events, err := repo.ListWithFilter(userFilter, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(user1Events), 2)

		// Filter by multiple criteria
		combinedFilter := EventFilter{
			Type:   &eventType,
			ToolID: &tool1ID,
			UserID: &user1ID,
		}
		filteredEvents, err := repo.ListWithFilter(combinedFilter, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(filteredEvents), 1)
		for _, event := range filteredEvents {
			assert.Equal(t, domain.EventTypeToolCheckedOut, event.Type)
			assert.Equal(t, &tool1ID, event.ToolID)
			assert.True(t,
				(event.UserID != nil && *event.UserID == user1ID) ||
					(event.ActorID != nil && *event.ActorID == user1ID),
			)
		}
	})

	t.Run("UUID Primary Keys", func(t *testing.T) {
		event1, err := repo.Create(domain.EventTypeToolCreated, &tool1ID, nil, &actorID, "Event 1", nil)
		require.NoError(t, err)

		event2, err := repo.Create(domain.EventTypeToolCreated, &tool2ID, nil, &actorID, "Event 2", nil)
		require.NoError(t, err)

		// UUIDs should be different
		assert.NotEqual(t, event1.ID, event2.ID)

		// UUIDs should be valid format (36 chars with dashes)
		assert.Len(t, event1.ID, 36)
		assert.Contains(t, event1.ID, "-")
	})
}

// TestPostgresEventRepo_ErrorCases tests error handling
func TestPostgresEventRepo_ErrorCases(t *testing.T) {
	db := setupSharedRepoTestDB(t)
	repo := NewPostgresEventRepo(db)

	t.Run("Get Non-existent Event", func(t *testing.T) {
		_, err := repo.Get("00000000-0000-0000-0000-000000000000")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event not found")
	})

	t.Run("Get Invalid UUID", func(t *testing.T) {
		_, err := repo.Get("invalid-uuid")
		assert.Error(t, err)
	})

	t.Run("Create Event with Invalid References", func(t *testing.T) {
		nonExistentID := "00000000-0000-0000-0000-000000000000"

		// This might fail due to foreign key constraints if they exist
		// The exact behavior depends on your database schema
		_, err := repo.Create(
			domain.EventTypeToolCheckedOut,
			&nonExistentID, // Non-existent tool
			&nonExistentID, // Non-existent user
			&nonExistentID, // Non-existent actor
			"This should fail",
			nil,
		)
		// We can't assert the exact error since it depends on FK constraints
		// but we can verify that some events work
		assert.NotNil(t, err) // Remove this line if FK constraints aren't enforced
	})

	t.Run("Empty Filter Returns All Events", func(t *testing.T) {
		// Create a test event
		userID := createTestUser(t, db, "Filter User", "filter@example.com", domain.UserRoleEmployee)
		_, err := repo.Create(domain.EventTypeUserCreated, nil, &userID, &userID, "Filter test", nil)
		require.NoError(t, err)

		// Empty filter should return events
		emptyFilter := EventFilter{}
		events, err := repo.ListWithFilter(emptyFilter, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(events), 1)
	})
}
