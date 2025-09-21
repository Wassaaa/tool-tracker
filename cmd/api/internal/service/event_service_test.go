package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/repo"
)

// TestEventService_CreateEvent tests the event creation workflow
func TestEventService_CreateEvent(t *testing.T) {
	t.Run("Successful event creation", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		userID := TestUserID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeToolCheckedOut, &toolID, &userID, &actorID, "Tool checked out for project")

		// Set expectations
		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolCheckedOut,
			&toolID,
			&userID,
			&actorID,
			"Tool checked out for project",
			(*string)(nil),
		).Return(createdEvent, nil)

		// Execute
		result, err := mocks.Service.CreateEvent(
			domain.EventTypeToolCheckedOut,
			&toolID,
			&userID,
			&actorID,
			"Tool checked out for project",
			nil,
		)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, createdEvent, result)
	})

	t.Run("Repository error should propagate", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		actorID := TestActorID
		repoError := assert.AnError
		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolCreated,
			&toolID,
			(*string)(nil),
			&actorID,
			"Tool created",
			(*string)(nil),
		).Return(domain.Event{}, repoError)

		_, err := mocks.Service.CreateEvent(
			domain.EventTypeToolCreated,
			&toolID,
			nil,
			&actorID,
			"Tool created",
			nil,
		)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})
}

// TestEventService_ListEvents tests the event listing with filters
func TestEventService_ListEvents(t *testing.T) {
	t.Run("List events without filters", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		userID := TestUserID
		expectedEvents := []domain.Event{
			CreateTestEvent("event1", domain.EventTypeToolCheckedOut, &toolID, &userID, nil, "Event 1"),
			CreateTestEvent("event2", domain.EventTypeToolCheckedIn, &toolID, &userID, nil, "Event 2"),
		}

		expectedFilter := repo.EventFilter{}
		mocks.MockRepo.EXPECT().ListWithFilter(expectedFilter, 50, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.ListEvents(50, 0, nil, nil, nil)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("List events with event type filter", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		userID := TestUserID
		eventTypeStr := string(domain.EventTypeToolCheckedOut)
		expectedEvents := []domain.Event{
			CreateTestEvent("event1", domain.EventTypeToolCheckedOut, &toolID, &userID, nil, "Event 1"),
		}

		eventTypeFilter := domain.EventTypeToolCheckedOut
		expectedFilter := repo.EventFilter{
			Type: &eventTypeFilter,
		}
		mocks.MockRepo.EXPECT().ListWithFilter(expectedFilter, 50, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.ListEvents(50, 0, &eventTypeStr, nil, nil)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("List events with tool ID filter", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		userID := TestUserID
		expectedEvents := []domain.Event{
			CreateTestEvent("event1", domain.EventTypeToolCheckedOut, &toolID, &userID, nil, "Event 1"),
		}

		expectedFilter := repo.EventFilter{
			ToolID: &toolID,
		}
		mocks.MockRepo.EXPECT().ListWithFilter(expectedFilter, 50, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.ListEvents(50, 0, nil, &toolID, nil)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("Default limit applied when zero", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		expectedFilter := repo.EventFilter{}
		mocks.MockRepo.EXPECT().ListWithFilter(expectedFilter, 50, 0).Return([]domain.Event{}, nil)

		_, err := mocks.Service.ListEvents(0, 0, nil, nil, nil)

		require.NoError(t, err)
	})

	t.Run("Maximum limit enforced", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		expectedFilter := repo.EventFilter{}
		mocks.MockRepo.EXPECT().ListWithFilter(expectedFilter, 500, 0).Return([]domain.Event{}, nil)

		_, err := mocks.Service.ListEvents(1000, 0, nil, nil, nil)

		require.NoError(t, err)
	})

	t.Run("Invalid event type should fail", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		invalidEventType := "invalid_event_type"

		_, err := mocks.Service.ListEvents(50, 0, &invalidEventType, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event type")
	})
}

// TestEventService_GetEvent tests event retrieval by ID
func TestEventService_GetEvent(t *testing.T) {
	t.Run("Successful event retrieval", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		expectedEvent := CreateTestEvent(TestEventID, domain.EventTypeToolCreated, nil, nil, nil, "Event retrieved")

		mocks.MockRepo.EXPECT().Get(TestEventID).Return(expectedEvent, nil)

		result, err := mocks.Service.GetEvent(TestEventID)

		require.NoError(t, err)
		assert.Equal(t, expectedEvent, result)
	})

	t.Run("Invalid event ID should fail", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.GetEvent(InvalidUUID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event_id must be a valid UUID")
	})
}

// TestEventService_GetToolHistory tests tool history retrieval
func TestEventService_GetToolHistory(t *testing.T) {
	t.Run("Successful tool history retrieval", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		expectedEvents := []domain.Event{
			CreateTestEvent("event1", domain.EventTypeToolCreated, &toolID, nil, nil, "Tool created"),
			CreateTestEvent("event2", domain.EventTypeToolCheckedOut, &toolID, nil, nil, "Tool checked out"),
		}

		mocks.MockRepo.EXPECT().ListByTool(TestToolID, 1000, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.GetToolHistory(TestToolID)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("Invalid tool ID should fail", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.GetToolHistory(InvalidUUID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool_id must be a valid UUID")
	})
}

// TestEventService_GetUserActivity tests user activity retrieval
func TestEventService_GetUserActivity(t *testing.T) {
	t.Run("Successful user activity retrieval", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		userID := TestUserID
		expectedEvents := []domain.Event{
			CreateTestEvent("event1", domain.EventTypeToolCheckedOut, nil, &userID, nil, "User checked out tool"),
			CreateTestEvent("event2", domain.EventTypeToolCheckedIn, nil, &userID, nil, "User checked in tool"),
		}

		mocks.MockRepo.EXPECT().ListByUser(TestUserID, 1000, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.GetUserActivity(TestUserID)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.GetUserActivity(InvalidUUID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}

// TestEventService_LogToolCreated tests the convenience logging method
func TestEventService_LogToolCreated(t *testing.T) {
	t.Run("Successful tool creation log", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeToolCreated, &toolID, nil, &actorID, "Tool created via API")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolCreated,
			&toolID,
			(*string)(nil),
			&actorID,
			"Tool created via API",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolCreated(TestToolID, TestActorID, "Tool created via API")

		require.NoError(t, err)
	})
}

// TestEventService_LogToolCheckedOut tests the checkout logging method
func TestEventService_LogToolCheckedOut(t *testing.T) {
	t.Run("Successful checkout log", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		userID := TestUserID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeToolCheckedOut, &toolID, &userID, &actorID, "Tool checked out for project")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolCheckedOut,
			&toolID,
			&userID,
			&actorID,
			"Tool checked out for project",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolCheckedOut(TestToolID, TestUserID, TestActorID, "Tool checked out for project")

		require.NoError(t, err)
	})
}

// TestEventService_GetEventsByType tests getting events by type
func TestEventService_GetEventsByType(t *testing.T) {
	t.Run("Successful get events by type", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		expectedEvents := []domain.Event{
			CreateTestEvent("event1", domain.EventTypeToolCreated, nil, nil, nil, "Tool created"),
			CreateTestEvent("event2", domain.EventTypeToolCreated, nil, nil, nil, "Another tool created"),
		}

		mocks.MockRepo.EXPECT().ListByType(domain.EventTypeToolCreated, 50, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.GetEventsByType(domain.EventTypeToolCreated, 50, 0)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("Default limit applied when zero", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().ListByType(domain.EventTypeToolCreated, 50, 0).Return([]domain.Event{}, nil)

		_, err := mocks.Service.GetEventsByType(domain.EventTypeToolCreated, 0, 0)

		require.NoError(t, err)
	})

	t.Run("Maximum limit enforced", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().ListByType(domain.EventTypeToolCreated, 500, 0).Return([]domain.Event{}, nil)

		_, err := mocks.Service.GetEventsByType(domain.EventTypeToolCreated, 600, 0)

		require.NoError(t, err)
	})

	t.Run("Negative offset corrected", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().ListByType(domain.EventTypeToolCreated, 50, 0).Return([]domain.Event{}, nil)

		_, err := mocks.Service.GetEventsByType(domain.EventTypeToolCreated, 50, -5)

		require.NoError(t, err)
	})
}

// TestEventService_GetEventCount tests getting event count
func TestEventService_GetEventCount(t *testing.T) {
	t.Run("Successful get event count", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		expectedCount := 100

		mocks.MockRepo.EXPECT().Count().Return(expectedCount, nil)

		result, err := mocks.Service.GetEventCount()

		require.NoError(t, err)
		assert.Equal(t, expectedCount, result)
	})

	t.Run("Repository error should propagate", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		repoError := assert.AnError
		mocks.MockRepo.EXPECT().Count().Return(0, repoError)

		_, err := mocks.Service.GetEventCount()

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})
}

// TestEventService_AdditionalLoggingMethods tests the remaining logging convenience methods
func TestEventService_AdditionalLoggingMethods(t *testing.T) {
	t.Run("LogToolUpdated", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeToolUpdated, &toolID, nil, &actorID, "Tool updated")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolUpdated,
			&toolID,
			(*string)(nil),
			&actorID,
			"Tool updated",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolUpdated(TestToolID, TestActorID, "Tool updated")

		require.NoError(t, err)
	})

	t.Run("LogToolDeleted", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeToolDeleted, &toolID, nil, &actorID, "Tool deleted")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolDeleted,
			&toolID,
			(*string)(nil),
			&actorID,
			"Tool deleted",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolDeleted(TestToolID, TestActorID, "Tool deleted")

		require.NoError(t, err)
	})

	t.Run("LogToolCheckedIn", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		userID := TestUserID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeToolCheckedIn, &toolID, &userID, &actorID, "Tool checked in")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolCheckedIn,
			&toolID,
			&userID,
			&actorID,
			"Tool checked in",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolCheckedIn(TestToolID, TestUserID, TestActorID, "Tool checked in")

		require.NoError(t, err)
	})

	t.Run("LogToolMaintenance", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		userID := TestUserID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeToolMaintenance, &toolID, &userID, nil, "Tool needs maintenance")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolMaintenance,
			&toolID,
			&userID,
			(*string)(nil),
			"Tool needs maintenance",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolMaintenance(TestToolID, TestUserID, "Tool needs maintenance")

		require.NoError(t, err)
	})

	t.Run("LogToolLost", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		toolID := TestToolID
		userID := TestUserID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeToolLost, &toolID, &userID, nil, "Tool lost")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolLost,
			&toolID,
			&userID,
			(*string)(nil),
			"Tool lost",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolLost(TestToolID, TestUserID, "Tool lost")

		require.NoError(t, err)
	})

	t.Run("LogUserCreated", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		userID := TestUserID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeUserCreated, nil, &userID, &actorID, "User created")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeUserCreated,
			(*string)(nil),
			&userID,
			&actorID,
			"User created",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogUserCreated(TestUserID, TestActorID, "User created")

		require.NoError(t, err)
	})

	t.Run("LogUserUpdated", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		userID := TestUserID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeUserUpdated, nil, &userID, &actorID, "User updated")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeUserUpdated,
			(*string)(nil),
			&userID,
			&actorID,
			"User updated",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogUserUpdated(TestUserID, TestActorID, "User updated")

		require.NoError(t, err)
	})

	t.Run("LogUserDeleted", func(t *testing.T) {
		mocks := SetupEventServiceMocks(t)
		defer mocks.Teardown()

		userID := TestUserID
		actorID := TestActorID
		createdEvent := CreateTestEvent(TestEventID, domain.EventTypeUserDeleted, nil, &userID, &actorID, "User deleted")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeUserDeleted,
			(*string)(nil),
			&userID,
			&actorID,
			"User deleted",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogUserDeleted(TestUserID, TestActorID, "User deleted")

		require.NoError(t, err)
	})
}
