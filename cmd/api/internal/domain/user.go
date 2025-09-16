package domain

import (
	"fmt"
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

	switch u.Role {
	case UserRoleEmployee, UserRoleAdmin, UserRoleManager:
		// Valid
	default:
		return fmt.Errorf("invalid role: %s", u.Role)
	}

	return nil
}
