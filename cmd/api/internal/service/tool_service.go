package service

import "github.com/wassaaa/tool-tracker/cmd/api/internal/domain"

type ToolRepo interface {
	Create(name string) (domain.Tool, error)
	List(limit, offset int) ([]domain.Tool, error)
}

type ToolService struct {
	Repo ToolRepo
}

func NewToolService(r ToolRepo) *ToolService {
	return &ToolService{Repo: r}
}

func (s *ToolService) CreateTool(name string) (domain.Tool, error) {
	return s.Repo.Create(name)
}

func (s *ToolService) ListTools(limit, offset int) ([]domain.Tool, error) {
	return s.Repo.List(limit, offset)
}
