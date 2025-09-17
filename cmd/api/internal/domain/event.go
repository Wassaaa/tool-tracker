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
	ToolID    *string   `json:"tool_id,omitempty"`
	UserID    *string   `json:"user_id,omitempty"`
	ActorID   *string   `json:"actor_id,omitempty"`
	Notes     string    `json:"notes"`
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// NewEvent constructs an Event and validates it.
func NewEvent(eventType EventType, toolID, userID, actorID *string, notes, metadata string) (Event, error) {
	e := Event{Type: eventType, ToolID: toolID, UserID: userID, ActorID: actorID, Notes: notes, Metadata: metadata}
	return e, e.Validate()
}

// ValidateEventType checks if the provided event type is valid.
func ValidateEventType(t EventType) error {
	if t == "" {
		return fmt.Errorf("%w: event type is required", ErrValidation)
	}
	if !t.IsValid() {
		return fmt.Errorf("%w: invalid event type %s", ErrValidation, t)
	}
	return nil
}

func (e *Event) Validate() error {
	if err := ValidateEventType(e.Type); err != nil {
		return err
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
