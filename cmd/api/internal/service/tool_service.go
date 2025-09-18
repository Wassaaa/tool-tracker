package service

import (
	"fmt"

	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type ToolRepo interface {
	Create(name string, status domain.ToolStatus) (domain.Tool, error)
	List(limit, offset int) ([]domain.Tool, error)
	Get(id string) (domain.Tool, error)
	Update(domain.Tool) (domain.Tool, error)
	Delete(id string) error
	ListByStatus(status domain.ToolStatus, limit, offset int) ([]domain.Tool, error)
	ListByUser(userID string, limit, offset int) ([]domain.Tool, error)
	Count() (int, error)
}

type ToolService struct {
	Repo   ToolRepo
	events EventLogger
}

// EventLogger provides event logging for tool lifecycle actions.
type EventLogger interface {
	LogToolCheckedOut(toolID string, userID string, actorID string, notes string) error
	LogToolCheckedIn(toolID string, userID string, actorID string, notes string) error
	LogToolMaintenance(toolID string, userID string, notes string) error
	LogToolLost(toolID string, userID string, notes string) error
	LogToolCreated(toolID string, actorID string, notes string) error
	LogToolUpdated(toolID string, actorID string, notes string) error
	LogToolDeleted(toolID string, actorID string, notes string) error
	LogUserCreated(userID string, actorID string, notes string) error
	LogUserUpdated(userID string, actorID string, notes string) error
	LogUserDeleted(userID string, actorID string, notes string) error
}

func NewToolService(r ToolRepo) *ToolService {
	return &ToolService{Repo: r}
}

// WithEventLogger sets the event logger dependency (optional chaining style).
func (s *ToolService) WithEventLogger(l EventLogger) *ToolService {
	s.events = l
	return s
}

func (s *ToolService) CreateTool(name string, status domain.ToolStatus, actorID, notes string) (domain.Tool, error) {
	t, err := domain.NewTool(name, status)
	if err != nil {
		return domain.Tool{}, err
	}
	created, err := s.Repo.Create(t.Name, t.Status)
	if err != nil {
		return domain.Tool{}, err
	}
	if s.events != nil {
		if created.ID != nil {
			_ = s.events.LogToolCreated(*created.ID, actorID, notes)
		}
	}
	return created, nil
}

func (s *ToolService) ListTools(limit, offset int) ([]domain.Tool, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.Repo.List(limit, offset)
}

func (s *ToolService) GetTool(id string) (domain.Tool, error) {
	if err := domain.ValidateUUID(id, "tool_id"); err != nil {
		return domain.Tool{}, err
	}
	return s.Repo.Get(id)
}

func (s *ToolService) UpdateTool(id string, name string, status domain.ToolStatus, actorID, notes string) (domain.Tool, error) {
	tool, err := s.applyAndSave(id, func(t *domain.Tool) error {
		t.Name = name
		t.Status = status
		return nil
	})
	if err != nil {
		return domain.Tool{}, err
	}
	if s.events != nil {
		if tool.ID != nil {
			_ = s.events.LogToolUpdated(*tool.ID, actorID, notes)
		}
	}
	return tool, nil
}

// CheckOutTool: internal controlled mutation (sets CurrentUserId, LastCheckedOutAt, Status)
func (s *ToolService) CheckOutTool(toolID, userID, actorID, notes string) (domain.Tool, error) {
	if err := domain.ValidateUUID(userID, "user_id"); err != nil {
		return domain.Tool{}, err
	}
	if actorID != "" && actorID != userID {
		if err := domain.ValidateUUID(actorID, "actor_id"); err != nil {
			return domain.Tool{}, err
		}
	}
	tool, err := s.applyAndSave(toolID, func(t *domain.Tool) error {
		if t.CurrentUserId != nil {
			return fmt.Errorf("%w: tool is already checked out", domain.ErrValidation)
		}
		if t.Status == domain.ToolStatusLost {
			return fmt.Errorf("%w: tool is marked as lost", domain.ErrValidation)
		}
		t.Status = domain.ToolStatusCheckedOut
		t.CurrentUserId = &userID
		return nil
	})
	if err != nil {
		return domain.Tool{}, err
	}
	if s.events != nil {
		_ = s.events.LogToolCheckedOut(toolID, userID, pickActor(actorID, userID), notes)
	}
	return tool, nil
}

// ReturnTool: clears checkout state
func (s *ToolService) ReturnTool(toolID, actorID, notes string) (domain.Tool, error) {
	var priorUserID string
	tool, err := s.applyAndSave(toolID, func(t *domain.Tool) error {
		if t.CurrentUserId == nil {
			return fmt.Errorf("%w: tool is already checked in", domain.ErrValidation)
		}
		// capture prior user id before clearing
		if t.CurrentUserId != nil {
			priorUserID = *t.CurrentUserId
		}
		t.CurrentUserId = nil
		t.Status = domain.ToolStatusInOffice
		return nil
	})
	if err != nil {
		return domain.Tool{}, err
	}
	if s.events != nil {
		_ = s.events.LogToolCheckedIn(toolID, priorUserID, pickActor(actorID, priorUserID), notes)
	}
	return tool, nil
}

// SendToMaintenance moves a tool to maintenance status.
func (s *ToolService) SendToMaintenance(toolID, actorID, notes string) (domain.Tool, error) {
	tool, err := s.applyAndSave(toolID, func(t *domain.Tool) error {
		if t.Status == domain.ToolStatusLost {
			return fmt.Errorf("%w: lost tools cannot be sent to maintenance", domain.ErrValidation)
		}
		if t.Status == domain.ToolStatusMaintenance {
			return nil
		}
		t.Status = domain.ToolStatusMaintenance
		return nil
	})
	if err != nil {
		return domain.Tool{}, err
	}
	if s.events != nil {
		_ = s.events.LogToolMaintenance(toolID, pickActor(actorID, ""), notes)
	}
	return tool, nil
}

// MarkLost marks a tool as lost.
func (s *ToolService) MarkLost(toolID, actorID, notes string) (domain.Tool, error) {
	tool, err := s.applyAndSave(toolID, func(t *domain.Tool) error {
		if t.Status == domain.ToolStatusLost {
			return nil
		}
		t.Status = domain.ToolStatusLost
		return nil
	})
	if err != nil {
		return domain.Tool{}, err
	}
	if s.events != nil {
		_ = s.events.LogToolLost(toolID, pickActor(actorID, ""), notes)
	}
	return tool, nil
}

// pickActor chooses actorID if provided, else fallback.
func pickActor(actorID, fallback string) string {
	if actorID != "" {
		return actorID
	}
	return fallback
}

func (s *ToolService) DeleteTool(id, actorID, notes string) error {
	if err := domain.ValidateUUID(id, "tool_id"); err != nil {
		return err
	}
	// load to get ID pointer value
	t, err := s.Repo.Get(id)
	if err != nil {
		return err
	}
	if err := s.Repo.Delete(id); err != nil {
		return err
	}
	if s.events != nil && t.ID != nil {
		_ = s.events.LogToolDeleted(*t.ID, actorID, notes)
	}
	return nil
}

func (s *ToolService) ListToolsByUser(userID string, limit, offset int) ([]domain.Tool, error) {
	if err := domain.ValidateUUID(userID, "user_id"); err != nil {
		return nil, err
	}
	return s.Repo.ListByUser(userID, limit, offset)
}

func (s *ToolService) ListToolsByStatus(status domain.ToolStatus, limit, offset int) ([]domain.Tool, error) {
	// Validate status
	if err := domain.ValidateToolStatus(status); err != nil {
		return nil, err
	}
	return s.Repo.ListByStatus(status, limit, offset)
}

func (s *ToolService) GetToolCount() (int, error) {
	return s.Repo.Count()
}

// applyAndSave centralizes: id validation, load, mutation, validation, timestamp, persist.
func (s *ToolService) applyAndSave(id string, mutate func(*domain.Tool) error) (domain.Tool, error) {
	if err := domain.ValidateUUID(id, "tool_id"); err != nil {
		return domain.Tool{}, err
	}
	current, err := s.Repo.Get(id)
	if err != nil {
		return domain.Tool{}, err
	}
	if err := mutate(&current); err != nil {
		return domain.Tool{}, err
	}
	if err := current.Validate(); err != nil {
		return domain.Tool{}, err
	}

	updated, err := s.Repo.Update(current)
	if err != nil {
		return domain.Tool{}, err
	}
	return updated, nil
}
