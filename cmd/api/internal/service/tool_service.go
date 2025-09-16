package service

import (
	"fmt"

	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type ToolRepo interface {
	Create(name string, status domain.ToolStatus) (domain.Tool, error)
	List(limit, offset int) ([]domain.Tool, error)
	Get(id string) (domain.Tool, error)
	Update(id string, name string, status domain.ToolStatus) (domain.Tool, error)
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
	// Validate input
	if name == "" {
		return domain.Tool{}, fmt.Errorf("tool name cannot be empty")
	}

	// Validate status
	tool := domain.Tool{Name: name, Status: status}
	if err := tool.Validate(); err != nil {
		return domain.Tool{}, err
	}

	return s.Repo.Create(name, status)
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
	if id == "" {
		return domain.Tool{}, fmt.Errorf("tool ID cannot be empty")
	}

	return s.Repo.Get(id)
}

func (s *ToolService) UpdateTool(id string, name string, status domain.ToolStatus) (domain.Tool, error) {
	if id == "" {
		return domain.Tool{}, fmt.Errorf("tool ID cannot be empty")
	}
	if name == "" {
		return domain.Tool{}, fmt.Errorf("tool name cannot be empty")
	}

	// Validate status
	tool := domain.Tool{Name: name, Status: status}
	if err := tool.Validate(); err != nil {
		return domain.Tool{}, err
	}

	return s.Repo.Update(id, name, status)
}

func (s *ToolService) DeleteTool(id string) error {
	if id == "" {
		return fmt.Errorf("tool ID cannot be empty")
	}

	return s.Repo.Delete(id)
}

func (s *ToolService) ListToolsByStatus(status domain.ToolStatus, limit, offset int) ([]domain.Tool, error) {
	// Validate status
	tool := domain.Tool{Status: status}
	if err := tool.Validate(); err != nil {
		return nil, err
	}

	return s.Repo.ListByStatus(status, limit, offset)
}

func (s *ToolService) GetToolCount() (int, error) {
	return s.Repo.Count()
}
