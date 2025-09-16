package service

import (
	"fmt"

	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/repo"
)

type EventRepo interface {
	Create(eventType domain.EventType, toolID *string, userID *string, actorID *string, notes string, metadata string) (domain.Event, error)
	List(limit, offset int) ([]domain.Event, error)
	Get(id string) (domain.Event, error)
	ListByType(eventType domain.EventType, limit, offset int) ([]domain.Event, error)
	ListByTool(toolID string, limit, offset int) ([]domain.Event, error)
	ListByUser(userID string, limit, offset int) ([]domain.Event, error)
	ListWithFilter(filter repo.EventFilter, limit, offset int) ([]domain.Event, error)
	Count() (int, error)
}

type EventService struct {
	Repo EventRepo
}

func NewEventService(r EventRepo) *EventService {
	return &EventService{Repo: r}
}

func (s *EventService) CreateEvent(eventType domain.EventType, toolID *string, userID *string, actorID *string, notes string, metadata string) (domain.Event, error) {
	// Validate event type
	if !eventType.IsValid() {
		return domain.Event{}, fmt.Errorf("invalid event type: %s", eventType)
	}

	// Validate using domain validation
	event := domain.Event{
		Type:     eventType,
		ToolID:   toolID,
		UserID:   userID,
		ActorID:  actorID,
		Notes:    notes,
		Metadata: metadata,
	}
	if err := event.Validate(); err != nil {
		return domain.Event{}, err
	}

	return s.Repo.Create(eventType, toolID, userID, actorID, notes, metadata)
}

func (s *EventService) ListEvents(limit, offset int, eventType *string, toolID *string, userID *string) ([]domain.Event, error) {
	if limit <= 0 {
		limit = 50 // Higher default for events
	}
	if limit > 500 {
		limit = 500 // Higher max for events
	}
	if offset < 0 {
		offset = 0
	}

	// Build filter
	filter := repo.EventFilter{}
	if eventType != nil && *eventType != "" {
		et := domain.EventType(*eventType)
		if !et.IsValid() {
			return nil, fmt.Errorf("invalid event type: %s", *eventType)
		}
		filter.Type = &et
	}
	if toolID != nil && *toolID != "" {
		filter.ToolID = toolID
	}
	if userID != nil && *userID != "" {
		filter.UserID = userID
	}

	return s.Repo.ListWithFilter(filter, limit, offset)
}

func (s *EventService) GetEvent(id string) (domain.Event, error) {
	if id == "" {
		return domain.Event{}, fmt.Errorf("event ID cannot be empty")
	}

	return s.Repo.Get(id)
}

func (s *EventService) GetToolHistory(toolID string) ([]domain.Event, error) {
	if toolID == "" {
		return nil, fmt.Errorf("tool ID cannot be empty")
	}

	// Get all events for this tool (higher limit for history)
	return s.Repo.ListByTool(toolID, 1000, 0)
}

func (s *EventService) GetUserActivity(userID string) ([]domain.Event, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Get all events for this user (higher limit for activity)
	return s.Repo.ListByUser(userID, 1000, 0)
}

func (s *EventService) GetEventsByType(eventType domain.EventType, limit, offset int) ([]domain.Event, error) {
	if !eventType.IsValid() {
		return nil, fmt.Errorf("invalid event type: %s", eventType)
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}

	return s.Repo.ListByType(eventType, limit, offset)
}

func (s *EventService) GetEventCount() (int, error) {
	return s.Repo.Count()
}

// Helper methods for specific event creation
func (s *EventService) LogToolCreated(toolID string, userID string) error {
	_, err := s.CreateEvent(domain.EventTypeToolCreated, &toolID, &userID, nil, "Tool created", "")
	return err
}

func (s *EventService) LogToolUpdated(toolID string, userID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolUpdated, &toolID, &userID, nil, notes, "")
	return err
}

func (s *EventService) LogToolDeleted(toolID string, userID string) error {
	_, err := s.CreateEvent(domain.EventTypeToolDeleted, &toolID, &userID, nil, "Tool deleted", "")
	return err
}

func (s *EventService) LogToolCheckedOut(toolID string, userID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolCheckedOut, &toolID, &userID, &actorID, notes, "")
	return err
}

func (s *EventService) LogToolCheckedIn(toolID string, userID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolCheckedIn, &toolID, &userID, &actorID, notes, "")
	return err
}

func (s *EventService) LogToolMaintenance(toolID string, userID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolMaintenance, &toolID, &userID, nil, notes, "")
	return err
}

func (s *EventService) LogToolLost(toolID string, userID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolLost, &toolID, &userID, nil, notes, "")
	return err
}

func (s *EventService) LogUserCreated(userID string, actorID string) error {
	_, err := s.CreateEvent(domain.EventTypeUserCreated, nil, &actorID, &userID, "User created", "")
	return err
}

func (s *EventService) LogUserUpdated(userID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeUserUpdated, nil, &actorID, &userID, notes, "")
	return err
}

func (s *EventService) LogUserDeleted(userID string, actorID string) error {
	_, err := s.CreateEvent(domain.EventTypeUserDeleted, nil, &actorID, &userID, "User deleted", "")
	return err
}
