package service

import (
	"fmt"

	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type UserRepo interface {
	Create(name string, email string, role domain.UserRole) (domain.User, error)
	List(limit, offset int) ([]domain.User, error)
	Get(id string) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
	Update(id string, name string, email string, role domain.UserRole) (domain.User, error)
	Delete(id string) error
	ListByRole(role domain.UserRole, limit, offset int) ([]domain.User, error)
	Count() (int, error)
}

type UserService struct {
	Repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{Repo: r}
}

func (s *UserService) CreateUser(name string, email string, role domain.UserRole) (domain.User, error) {
	u, err := domain.NewUser(name, email, role)
	if err != nil {
		return domain.User{}, err
	}

	// Check if user with email already exists
	if existing, err := s.Repo.GetByEmail(u.Email); err == nil && existing.ID != "" {
		return domain.User{}, fmt.Errorf("%w: user with email '%s' already exists", domain.ErrConflict, u.Email)
	}

	return s.Repo.Create(u.Name, u.Email, u.Role)
}

func (s *UserService) ListUsers(limit, offset int) ([]domain.User, error) {
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

func (s *UserService) GetUser(id string) (domain.User, error) {
	if err := domain.ValidateUUID(id, "user_id"); err != nil {
		return domain.User{}, err
	}
	return s.Repo.Get(id)
}

func (s *UserService) GetUserByEmail(email string) (domain.User, error) {
	if email == "" {
		return domain.User{}, fmt.Errorf("user email cannot be empty")
	}

	return s.Repo.GetByEmail(email)
}

func (s *UserService) UpdateUser(id string, name string, email string, role domain.UserRole) (domain.User, error) {
	if err := domain.ValidateUUID(id, "user_id"); err != nil {
		return domain.User{}, err
	}
	u, err := domain.NewUser(name, email, role)
	if err != nil {
		return domain.User{}, err
	}

	// Ensure email uniqueness (if another user has this email)
	if existing, err := s.Repo.GetByEmail(u.Email); err == nil && existing.ID != id {
		return domain.User{}, fmt.Errorf("%w: user with email '%s' already exists", domain.ErrConflict, u.Email)
	}

	return s.Repo.Update(id, u.Name, u.Email, u.Role)
}

func (s *UserService) DeleteUser(id string) error {
	if err := domain.ValidateUUID(id, "user_id"); err != nil {
		return err
	}
	return s.Repo.Delete(id)
}

func (s *UserService) ListUsersByRole(role domain.UserRole, limit, offset int) ([]domain.User, error) {
	// Validate role
	if err := domain.ValidateUserRole(role); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.Repo.ListByRole(role, limit, offset)
}

func (s *UserService) GetUserCount() (int, error) {
	return s.Repo.Count()
}
