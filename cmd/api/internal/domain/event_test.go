package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEventType_IsValid tests the EventType validation
func TestEventType_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		expected  bool
	}{
		// Tool events
		{"Valid TOOL_CREATED", EventTypeToolCreated, true},
		{"Valid TOOL_UPDATED", EventTypeToolUpdated, true},
		{"Valid TOOL_DELETED", EventTypeToolDeleted, true},
		{"Valid TOOL_CHECKED_OUT", EventTypeToolCheckedOut, true},
		{"Valid TOOL_CHECKED_IN", EventTypeToolCheckedIn, true},
		{"Valid TOOL_MAINTENANCE", EventTypeToolMaintenance, true},
		{"Valid TOOL_LOST", EventTypeToolLost, true},
		// User events
		{"Valid USER_CREATED", EventTypeUserCreated, true},
		{"Valid USER_UPDATED", EventTypeUserUpdated, true},
		{"Valid USER_DELETED", EventTypeUserDeleted, true},
		// Invalid cases
		{"Invalid empty", EventType(""), false},
		{"Invalid random", EventType("INVALID_EVENT"), false},
		{"Invalid case sensitive", EventType("tool_created"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.eventType.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNewEvent tests event creation
func TestNewEvent(t *testing.T) {
	toolID := "tool-123"
	userID := "user-456"
	actorID := "actor-789"
	metadata := `{"key": "value"}`

	t.Run("Valid event creation with all fields", func(t *testing.T) {
		event, err := NewEvent(
			EventTypeToolCheckedOut,
			&toolID,
			&userID,
			&actorID,
			"Tool checked out to user",
			&metadata,
		)

		require.NoError(t, err)
		assert.Equal(t, EventTypeToolCheckedOut, event.Type)
		assert.Equal(t, &toolID, event.ToolID)
		assert.Equal(t, &userID, event.UserID)
		assert.Equal(t, &actorID, event.ActorID)
		assert.Equal(t, "Tool checked out to user", event.Notes)
		assert.Equal(t, &metadata, event.Metadata)
	})

	t.Run("Valid event creation with minimal fields", func(t *testing.T) {
		event, err := NewEvent(
			EventTypeUserCreated,
			nil,
			&userID,
			nil,
			"User created",
			nil,
		)

		require.NoError(t, err)
		assert.Equal(t, EventTypeUserCreated, event.Type)
		assert.Nil(t, event.ToolID)
		assert.Equal(t, &userID, event.UserID)
		assert.Nil(t, event.ActorID)
		assert.Equal(t, "User created", event.Notes)
		assert.Nil(t, event.Metadata)
	})

	t.Run("Invalid event type should fail", func(t *testing.T) {
		_, err := NewEvent(
			EventType("INVALID"),
			&toolID,
			&userID,
			&actorID,
			"Some notes",
			nil,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid event type")
	})

	t.Run("Empty event type should fail", func(t *testing.T) {
		_, err := NewEvent(
			EventType(""),
			&toolID,
			&userID,
			&actorID,
			"Some notes",
			nil,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event type is required")
	})
}

// TestEvent_Validate tests the Event Validate method comprehensively
func TestEvent_Validate(t *testing.T) {
	t.Run("Valid event should pass validation", func(t *testing.T) {
		event := Event{
			Type:  EventTypeToolCreated,
			Notes: "Tool was created",
		}

		err := event.Validate()
		assert.NoError(t, err)
	})

	t.Run("Empty event type should fail validation", func(t *testing.T) {
		event := Event{
			Type:  EventType(""),
			Notes: "Some notes",
		}

		err := event.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event type is required")
	})

	t.Run("All valid event types should pass validation", func(t *testing.T) {
		validTypes := ValidEventTypes()

		for _, eventType := range validTypes {
			event := Event{
				Type:  eventType,
				Notes: "Test notes",
			}

			err := event.Validate()
			assert.NoError(t, err, "Event type %s should be valid", eventType)
		}
	})

	t.Run("Invalid event type should fail validation", func(t *testing.T) {
		invalidTypes := []EventType{
			EventType(""),
			EventType("INVALID"),
			EventType("tool_created"), // case sensitive
			EventType("RANDOM_EVENT"),
		}

		for _, eventType := range invalidTypes {
			event := Event{
				Type:  eventType,
				Notes: "Test notes",
			}

			err := event.Validate()
			assert.Error(t, err, "Event type %s should be invalid", eventType)
			if eventType == EventType("") {
				assert.Contains(t, err.Error(), "event type is required")
			} else {
				assert.Contains(t, err.Error(), "invalid event type")
			}
		}
	})
}

// TestValidEventTypes tests the helper function
func TestValidEventTypes(t *testing.T) {
	types := ValidEventTypes()

	assert.Len(t, types, 10)

	// Check tool events
	assert.Contains(t, types, EventTypeToolCreated)
	assert.Contains(t, types, EventTypeToolUpdated)
	assert.Contains(t, types, EventTypeToolDeleted)
	assert.Contains(t, types, EventTypeToolCheckedOut)
	assert.Contains(t, types, EventTypeToolCheckedIn)
	assert.Contains(t, types, EventTypeToolMaintenance)
	assert.Contains(t, types, EventTypeToolLost)

	// Check user events
	assert.Contains(t, types, EventTypeUserCreated)
	assert.Contains(t, types, EventTypeUserUpdated)
	assert.Contains(t, types, EventTypeUserDeleted)
}

// TestEventTypeCategories tests logical groupings
func TestEventTypeCategories(t *testing.T) {
	toolEvents := []EventType{
		EventTypeToolCreated,
		EventTypeToolUpdated,
		EventTypeToolDeleted,
		EventTypeToolCheckedOut,
		EventTypeToolCheckedIn,
		EventTypeToolMaintenance,
		EventTypeToolLost,
	}

	userEvents := []EventType{
		EventTypeUserCreated,
		EventTypeUserUpdated,
		EventTypeUserDeleted,
	}

	t.Run("All tool events should be valid", func(t *testing.T) {
		for _, eventType := range toolEvents {
			assert.True(t, eventType.IsValid(), "Tool event %s should be valid", eventType)
		}
	})

	t.Run("All user events should be valid", func(t *testing.T) {
		for _, eventType := range userEvents {
			assert.True(t, eventType.IsValid(), "User event %s should be valid", eventType)
		}
	})
}
