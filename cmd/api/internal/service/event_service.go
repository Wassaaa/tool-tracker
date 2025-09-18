package service

import (
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/repo"
)

type EventRepo interface {
	Create(eventType domain.EventType, toolID *string, userID *string, actorID *string, notes string, metadata *string) (domain.Event, error)
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

func (s *EventService) CreateEvent(eventType domain.EventType, toolID *string, userID *string, actorID *string, notes string, metadata *string) (domain.Event, error) {
	evt, err := domain.NewEvent(eventType, toolID, userID, actorID, notes, metadata)
	if err != nil {
		return domain.Event{}, err
	}
	return s.Repo.Create(evt.Type, evt.ToolID, evt.UserID, evt.ActorID, evt.Notes, evt.Metadata)
}

func (s *EventService) ListEvents(limit, offset int, eventType *string, toolID *string, userID *string) ([]domain.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}

	// Build filter
	filter := repo.EventFilter{}
	if eventType != nil && *eventType != "" {
		et := domain.EventType(*eventType)
		if err := domain.ValidateEventType(et); err != nil {
			return nil, err
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
	if err := domain.ValidateUUID(id, "event_id"); err != nil {
		return domain.Event{}, err
	}
	return s.Repo.Get(id)
}

func (s *EventService) GetToolHistory(toolID string) ([]domain.Event, error) {
	if err := domain.ValidateUUID(toolID, "tool_id"); err != nil {
		return nil, err
	}

	// Get all events for this tool (higher limit for history)
	return s.Repo.ListByTool(toolID, 1000, 0)
}

func (s *EventService) GetUserActivity(userID string) ([]domain.Event, error) {
	if err := domain.ValidateUUID(userID, "user_id"); err != nil {
		return nil, err
	}

	// Get all events for this user (higher limit for activity)
	return s.Repo.ListByUser(userID, 1000, 0)
}

func (s *EventService) GetEventsByType(eventType domain.EventType, limit, offset int) ([]domain.Event, error) {
	if err := domain.ValidateEventType(eventType); err != nil {
		return nil, err
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
// Tool CRUD logs (actor-aware)
func (s *EventService) LogToolCreated(toolID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolCreated, &toolID, nil, &actorID, notes, nil)
	return err
}

func (s *EventService) LogToolUpdated(toolID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolUpdated, &toolID, nil, &actorID, notes, nil)
	return err
}

func (s *EventService) LogToolDeleted(toolID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolDeleted, &toolID, nil, &actorID, notes, nil)
	return err
}

func (s *EventService) LogToolCheckedOut(toolID string, userID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolCheckedOut, &toolID, &userID, &actorID, notes, nil)
	return err
}

// Tool action logs
func (s *EventService) LogToolCheckedIn(toolID string, userID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolCheckedIn, &toolID, &userID, &actorID, notes, nil)
	return err
}

func (s *EventService) LogToolMaintenance(toolID string, userID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolMaintenance, &toolID, &userID, nil, notes, nil)
	return err
}

func (s *EventService) LogToolLost(toolID string, userID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeToolLost, &toolID, &userID, nil, notes, nil)
	return err
}

// User CRUD logs
func (s *EventService) LogUserCreated(userID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeUserCreated, nil, &userID, &actorID, notes, nil)
	return err
}

func (s *EventService) LogUserUpdated(userID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeUserUpdated, nil, &userID, &actorID, notes, nil)
	return err
}

func (s *EventService) LogUserDeleted(userID string, actorID string, notes string) error {
	_, err := s.CreateEvent(domain.EventTypeUserDeleted, nil, &userID, &actorID, notes, nil)
	return err
}
