package domain

import (
	"fmt"
	"time"
)

type EventType string

const (
	EventTypeToolCreated     EventType = "TOOL_CREATED"
	EventTypeToolUpdated     EventType = "TOOL_UPDATED"
	EventTypeToolDeleted     EventType = "TOOL_DELETED"
	EventTypeToolCheckedOut  EventType = "TOOL_CHECKED_OUT"
	EventTypeToolCheckedIn   EventType = "TOOL_CHECKED_IN"
	EventTypeToolMaintenance EventType = "TOOL_MAINTENANCE"
	EventTypeToolLost        EventType = "TOOL_LOST"
	EventTypeUserCreated     EventType = "USER_CREATED"
	EventTypeUserUpdated     EventType = "USER_UPDATED"
	EventTypeUserDeleted     EventType = "USER_DELETED"
)

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	ToolID    *string   `json:"tool_id,omitempty"`  // Pointer for optional field
	UserID    *string   `json:"user_id,omitempty"`  // Who performed the action
	ActorID   *string   `json:"actor_id,omitempty"` // Who the action was performed on
	Notes     string    `json:"notes"`
	Metadata  string    `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt time.Time `json:"created_at"`
}

func (e *Event) Validate() error {
	if e.Type == "" {
		return fmt.Errorf("event type is required")
	}

	switch e.Type {
	case EventTypeToolCreated, EventTypeToolUpdated, EventTypeToolDeleted,
		EventTypeToolCheckedOut, EventTypeToolCheckedIn, EventTypeToolMaintenance, EventTypeToolLost,
		EventTypeUserCreated, EventTypeUserUpdated, EventTypeUserDeleted:
		// Valid
	default:
		return fmt.Errorf("invalid event type: %s", e.Type)
	}

	return nil
}

func (t EventType) IsValid() bool {
	switch t {
	case EventTypeToolCreated, EventTypeToolUpdated, EventTypeToolDeleted,
		EventTypeToolCheckedOut, EventTypeToolCheckedIn, EventTypeToolMaintenance, EventTypeToolLost,
		EventTypeUserCreated, EventTypeUserUpdated, EventTypeUserDeleted:
		return true
	default:
		return false
	}
}

// Helper function to get all valid event types
func ValidEventTypes() []EventType {
	return []EventType{
		EventTypeToolCreated,
		EventTypeToolUpdated,
		EventTypeToolDeleted,
		EventTypeToolCheckedOut,
		EventTypeToolCheckedIn,
		EventTypeToolMaintenance,
		EventTypeToolLost,
		EventTypeUserCreated,
		EventTypeUserUpdated,
		EventTypeUserDeleted,
	}
}
