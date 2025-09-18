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
	Count() (int, error)
}

type ToolService struct {
	Repo ToolRepo
}

func NewToolService(r ToolRepo) *ToolService {
	return &ToolService{Repo: r}
}

func (s *ToolService) CreateTool(name string, status domain.ToolStatus) (domain.Tool, error) {
	t, err := domain.NewTool(name, status)
	if err != nil {
		return domain.Tool{}, err
	}

	return s.Repo.Create(t.Name, t.Status)
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

func (s *ToolService) UpdateTool(id string, name string, status domain.ToolStatus) (domain.Tool, error) {
	return s.applyAndSave(id, func(t *domain.Tool) error {
		t.Name = name
		t.Status = status
		return nil
	})
}

// CheckOutTool: internal controlled mutation (sets CurrentUserId, LastCheckedOutAt, Status)
func (s *ToolService) CheckOutTool(toolID, userID string) (domain.Tool, error) {
	if err := domain.ValidateUUID(userID, "user_id"); err != nil {
		return domain.Tool{}, err
	}
	return s.applyAndSave(toolID, func(t *domain.Tool) error {
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
}

// ReturnTool: clears checkout state
func (s *ToolService) ReturnTool(toolID string) (domain.Tool, error) {
	return s.applyAndSave(toolID, func(t *domain.Tool) error {
		if t.CurrentUserId == nil {
			return fmt.Errorf("%w: tool is already checked in", domain.ErrValidation)
		}
		t.CurrentUserId = nil
		t.Status = domain.ToolStatusInOffice
		return nil
	})
}

// SendToMaintenance moves a tool to maintenance status.
func (s *ToolService) SendToMaintenance(toolID string) (domain.Tool, error) {
	return s.applyAndSave(toolID, func(t *domain.Tool) error {
		if t.Status == domain.ToolStatusLost {
			return fmt.Errorf("%w: lost tools cannot be sent to maintenance", domain.ErrValidation)
		}
		if t.Status == domain.ToolStatusMaintenance {
			return nil
		}
		// Not clearing current_user_id here
		t.Status = domain.ToolStatusMaintenance
		return nil
	})
}

// MarkLost marks a tool as lost.
func (s *ToolService) MarkLost(toolID string) (domain.Tool, error) {
	return s.applyAndSave(toolID, func(t *domain.Tool) error {
		if t.Status == domain.ToolStatusLost {
			return nil
		}
		t.Status = domain.ToolStatusLost
		return nil
	})
}

func (s *ToolService) DeleteTool(id string) error {
	if err := domain.ValidateUUID(id, "tool_id"); err != nil {
		return err
	}
	return s.Repo.Delete(id)
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
