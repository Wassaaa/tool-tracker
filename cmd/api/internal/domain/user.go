package domain

import (
	"fmt"
	"regexp"
	"time"
)

type UserRole string

const (
	UserRoleEmployee UserRole = "EMPLOYEE"
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleManager  UserRole = "MANAGER"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser constructs a User with defaults and validates it.
func NewUser(name, email string, role UserRole) (User, error) {
	if role == "" {
		role = UserRoleEmployee
	}
	u := User{Name: name, Email: email, Role: role}
	return u, u.Validate()
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}

	if u.Email == "" {
		return fmt.Errorf("%w: email is required", ErrValidation)
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("%w: invalid email format", ErrValidation)
	}

	if err := ValidateUserRole(u.Role); err != nil {
		return err
	}

	return nil
}

func ValidateUserRole(r UserRole) error {
	if !r.IsValid() {
		return fmt.Errorf("%w: invalid role %s", ErrValidation, r)
	}

	return nil
}

func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleEmployee, UserRoleAdmin, UserRoleManager:
		return true
	default:
		return false
	}
}

// Helper function to get all valid roles
func ValidUserRoles() []UserRole {
	return []UserRole{
		UserRoleEmployee,
		UserRoleAdmin,
		UserRoleManager,
	}
}
