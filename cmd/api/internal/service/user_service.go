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
	// Validate input
	if name == "" {
		return domain.User{}, fmt.Errorf("user name cannot be empty")
	}
	if email == "" {
		return domain.User{}, fmt.Errorf("user email cannot be empty")
	}

	// Set default role if empty
	if role == "" {
		role = domain.UserRoleEmployee
	}

	// Validate using domain validation
	user := domain.User{Name: name, Email: email, Role: role}
	if err := user.Validate(); err != nil {
		return domain.User{}, err
	}

	// Check if user with email already exists
	_, err := s.Repo.GetByEmail(email)
	if err == nil {
		return domain.User{}, fmt.Errorf("user with email '%s' already exists", email)
	}

	return s.Repo.Create(name, email, role)
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
	if id == "" {
		return domain.User{}, fmt.Errorf("user ID cannot be empty")
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
	if id == "" {
		return domain.User{}, fmt.Errorf("user ID cannot be empty")
	}
	if name == "" {
		return domain.User{}, fmt.Errorf("user name cannot be empty")
	}
	if email == "" {
		return domain.User{}, fmt.Errorf("user email cannot be empty")
	}

	// Validate using domain validation
	user := domain.User{Name: name, Email: email, Role: role}
	if err := user.Validate(); err != nil {
		return domain.User{}, err
	}

	// Check if trying to update to an email that already exists (for different user)
	existingUser, err := s.Repo.GetByEmail(email)
	if err == nil && existingUser.ID != id {
		return domain.User{}, fmt.Errorf("user with email '%s' already exists", email)
	}

	return s.Repo.Update(id, name, email, role)
}

func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	return s.Repo.Delete(id)
}

func (s *UserService) ListUsersByRole(role domain.UserRole, limit, offset int) ([]domain.User, error) {
	// Validate role
	if !role.IsValid() {
		return nil, fmt.Errorf("invalid role: %s", role)
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
