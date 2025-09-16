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

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}

	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("invalid email format")
	}

	switch u.Role {
	case UserRoleEmployee, UserRoleAdmin, UserRoleManager:
		// Valid
	default:
		return fmt.Errorf("invalid role: %s", u.Role)
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
