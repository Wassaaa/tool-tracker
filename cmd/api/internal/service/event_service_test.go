package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/repo"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/service/mocks"
)

// Helper function to create a valid event for testing
func createTestEvent(id string, eventType domain.EventType, toolID, userID, actorID *string, notes string) domain.Event {
	now := time.Now()
	return domain.Event{
		ID:        id,
		Type:      eventType,
		ToolID:    toolID,
		UserID:    userID,
		ActorID:   actorID,
		Notes:     notes,
		Metadata:  nil,
		CreatedAt: now,
	}
}

// EventTestMocks holds all the mock dependencies for event service testing
type EventTestMocks struct {
	Ctrl     *gomock.Controller
	MockRepo *mocks.MockEventRepo
	Service  *EventService
}

// setupEventMocks creates all necessary mocks for event service testing
func setupEventMocks(t *testing.T) *EventTestMocks {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockEventRepo(ctrl)

	return &EventTestMocks{
		Ctrl:     ctrl,
		MockRepo: mockRepo,
		Service:  NewEventService(mockRepo),
	}
}

// teardown cleans up the mocks
func (etm *EventTestMocks) teardown() {
	etm.Ctrl.Finish()
}

// TestEventService_CreateEvent tests the event creation workflow
func TestEventService_CreateEvent(t *testing.T) {
	toolID := "123e4567-e89b-12d3-a456-426614174000"
	userID := "456e7890-e89b-12d3-a456-426614174000"
	actorID := "789e0123-e89b-12d3-a456-426614174000"

	t.Run("Successful event creation", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		eventID := "abc12345-e89b-12d3-a456-426614174000"
		createdEvent := createTestEvent(eventID, domain.EventTypeToolCheckedOut, &toolID, &userID, &actorID, "Tool checked out for project")

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
		mocks := setupEventMocks(t)
		defer mocks.teardown()

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
	toolID := "123e4567-e89b-12d3-a456-426614174000"
	userID := "456e7890-e89b-12d3-a456-426614174000"

	t.Run("List events without filters", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		expectedEvents := []domain.Event{
			createTestEvent("event1", domain.EventTypeToolCheckedOut, &toolID, &userID, nil, "Event 1"),
			createTestEvent("event2", domain.EventTypeToolCheckedIn, &toolID, &userID, nil, "Event 2"),
		}

		expectedFilter := repo.EventFilter{}
		mocks.MockRepo.EXPECT().ListWithFilter(expectedFilter, 50, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.ListEvents(50, 0, nil, nil, nil)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("List events with event type filter", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		eventTypeStr := string(domain.EventTypeToolCheckedOut)
		expectedEvents := []domain.Event{
			createTestEvent("event1", domain.EventTypeToolCheckedOut, &toolID, &userID, nil, "Event 1"),
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
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		expectedEvents := []domain.Event{
			createTestEvent("event1", domain.EventTypeToolCheckedOut, &toolID, &userID, nil, "Event 1"),
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
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		expectedFilter := repo.EventFilter{}
		mocks.MockRepo.EXPECT().ListWithFilter(expectedFilter, 50, 0).Return([]domain.Event{}, nil)

		_, err := mocks.Service.ListEvents(0, 0, nil, nil, nil)

		require.NoError(t, err)
	})

	t.Run("Maximum limit enforced", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		expectedFilter := repo.EventFilter{}
		mocks.MockRepo.EXPECT().ListWithFilter(expectedFilter, 500, 0).Return([]domain.Event{}, nil)

		_, err := mocks.Service.ListEvents(1000, 0, nil, nil, nil)

		require.NoError(t, err)
	})

	t.Run("Invalid event type should fail", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		invalidEventType := "invalid_event_type"

		_, err := mocks.Service.ListEvents(50, 0, &invalidEventType, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event type")
	})
}

// TestEventService_GetEvent tests event retrieval by ID
func TestEventService_GetEvent(t *testing.T) {
	eventID := "123e4567-e89b-12d3-a456-426614174000"

	t.Run("Successful event retrieval", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		expectedEvent := createTestEvent(eventID, domain.EventTypeToolCreated, nil, nil, nil, "Event retrieved")

		mocks.MockRepo.EXPECT().Get(eventID).Return(expectedEvent, nil)

		result, err := mocks.Service.GetEvent(eventID)

		require.NoError(t, err)
		assert.Equal(t, expectedEvent, result)
	})

	t.Run("Invalid event ID should fail", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		_, err := mocks.Service.GetEvent("invalid-uuid")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event_id must be a valid UUID")
	})
}

// TestEventService_GetToolHistory tests tool history retrieval
func TestEventService_GetToolHistory(t *testing.T) {
	toolID := "123e4567-e89b-12d3-a456-426614174000"

	t.Run("Successful tool history retrieval", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		expectedEvents := []domain.Event{
			createTestEvent("event1", domain.EventTypeToolCreated, &toolID, nil, nil, "Tool created"),
			createTestEvent("event2", domain.EventTypeToolCheckedOut, &toolID, nil, nil, "Tool checked out"),
		}

		mocks.MockRepo.EXPECT().ListByTool(toolID, 1000, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.GetToolHistory(toolID)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("Invalid tool ID should fail", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		_, err := mocks.Service.GetToolHistory("invalid-uuid")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool_id must be a valid UUID")
	})
}

// TestEventService_GetUserActivity tests user activity retrieval
func TestEventService_GetUserActivity(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"

	t.Run("Successful user activity retrieval", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		expectedEvents := []domain.Event{
			createTestEvent("event1", domain.EventTypeToolCheckedOut, nil, &userID, nil, "User checked out tool"),
			createTestEvent("event2", domain.EventTypeToolCheckedIn, nil, &userID, nil, "User checked in tool"),
		}

		mocks.MockRepo.EXPECT().ListByUser(userID, 1000, 0).Return(expectedEvents, nil)

		result, err := mocks.Service.GetUserActivity(userID)

		require.NoError(t, err)
		assert.Equal(t, expectedEvents, result)
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		_, err := mocks.Service.GetUserActivity("invalid-uuid")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}

// TestEventService_LogToolCreated tests the convenience logging method
func TestEventService_LogToolCreated(t *testing.T) {
	toolID := "123e4567-e89b-12d3-a456-426614174000"
	actorID := "789e0123-e89b-12d3-a456-426614174000"

	t.Run("Successful tool creation log", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		eventID := "abc12345-e89b-12d3-a456-426614174000"
		createdEvent := createTestEvent(eventID, domain.EventTypeToolCreated, &toolID, nil, &actorID, "Tool created via API")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolCreated,
			&toolID,
			(*string)(nil),
			&actorID,
			"Tool created via API",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolCreated(toolID, actorID, "Tool created via API")

		require.NoError(t, err)
	})
}

// TestEventService_LogToolCheckedOut tests the checkout logging method
func TestEventService_LogToolCheckedOut(t *testing.T) {
	toolID := "123e4567-e89b-12d3-a456-426614174000"
	userID := "456e7890-e89b-12d3-a456-426614174000"
	actorID := "789e0123-e89b-12d3-a456-426614174000"

	t.Run("Successful checkout log", func(t *testing.T) {
		mocks := setupEventMocks(t)
		defer mocks.teardown()

		eventID := "abc12345-e89b-12d3-a456-426614174000"
		createdEvent := createTestEvent(eventID, domain.EventTypeToolCheckedOut, &toolID, &userID, &actorID, "Tool checked out for project")

		mocks.MockRepo.EXPECT().Create(
			domain.EventTypeToolCheckedOut,
			&toolID,
			&userID,
			&actorID,
			"Tool checked out for project",
			(*string)(nil),
		).Return(createdEvent, nil)

		err := mocks.Service.LogToolCheckedOut(toolID, userID, actorID, "Tool checked out for project")

		require.NoError(t, err)
	})
}
