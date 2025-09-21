package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserRole_IsValid tests the UserRole validation
func TestUserRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"Valid EMPLOYEE", UserRoleEmployee, true},
		{"Valid ADMIN", UserRoleAdmin, true},
		{"Valid MANAGER", UserRoleManager, true},
		{"Invalid empty", UserRole(""), false},
		{"Invalid random", UserRole("INVALID_ROLE"), false},
		{"Invalid case sensitive", UserRole("employee"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNewUser tests user creation
func TestNewUser(t *testing.T) {
	t.Run("Valid user creation", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", UserRoleEmployee)

		require.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, UserRoleEmployee, user.Role)
	})

	t.Run("Default role when empty", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "")

		require.NoError(t, err)
		assert.Equal(t, UserRoleEmployee, user.Role)
	})

	t.Run("Empty name should fail", func(t *testing.T) {
		_, err := NewUser("", "john@example.com", UserRoleEmployee)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("Empty email should fail", func(t *testing.T) {
		_, err := NewUser("John Doe", "", UserRoleEmployee)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("Invalid email format should fail", func(t *testing.T) {
		invalidEmails := []string{
			"invalid-email",
			"@example.com",
			"john@",
			"john.example.com",
			"john@example",
			"john doe@example.com",
		}

		for _, email := range invalidEmails {
			_, err := NewUser("John Doe", email, UserRoleEmployee)

			assert.Error(t, err, "Email %s should be invalid", email)
			assert.Contains(t, err.Error(), "invalid email format")
		}
	})

	t.Run("Valid email formats should pass", func(t *testing.T) {
		validEmails := []string{
			"john@example.com",
			"john.doe@example.com",
			"john+tag@example.co.uk",
			"john123@example123.com",
			"j@e.co",
		}

		for _, email := range validEmails {
			user, err := NewUser("John Doe", email, UserRoleEmployee)

			require.NoError(t, err, "Email %s should be valid", email)
			assert.Equal(t, email, user.Email)
		}
	})

	t.Run("Invalid role should fail", func(t *testing.T) {
		_, err := NewUser("John Doe", "john@example.com", UserRole("INVALID"))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
	})
}

// TestUser_Validate tests the User Validate method comprehensively
func TestUser_Validate(t *testing.T) {
	t.Run("Valid user should pass validation", func(t *testing.T) {
		user := User{
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  UserRoleEmployee,
		}

		err := user.Validate()
		assert.NoError(t, err)
	})

	t.Run("Empty name should fail validation", func(t *testing.T) {
		user := User{
			Name:  "",
			Email: "john@example.com",
			Role:  UserRoleEmployee,
		}

		err := user.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("Empty email should fail validation", func(t *testing.T) {
		user := User{
			Name:  "John Doe",
			Email: "",
			Role:  UserRoleEmployee,
		}

		err := user.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("Invalid email should fail validation", func(t *testing.T) {
		user := User{
			Name:  "John Doe",
			Email: "invalid-email",
			Role:  UserRoleEmployee,
		}

		err := user.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("All valid roles should pass validation", func(t *testing.T) {
		validRoles := []UserRole{
			UserRoleEmployee,
			UserRoleAdmin,
			UserRoleManager,
		}

		for _, role := range validRoles {
			user := User{
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  role,
			}

			err := user.Validate()
			assert.NoError(t, err, "Role %s should be valid", role)
		}
	})

	t.Run("Invalid role should fail validation", func(t *testing.T) {
		invalidRoles := []UserRole{
			UserRole(""),
			UserRole("INVALID"),
			UserRole("employee"), // case sensitive
			UserRole("SUPER_ADMIN"),
		}

		for _, role := range invalidRoles {
			user := User{
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  role,
			}

			err := user.Validate()
			assert.Error(t, err, "Role %s should be invalid", role)
			assert.Contains(t, err.Error(), "invalid role")
		}
	})

	t.Run("Multiple validation errors - should fail on first", func(t *testing.T) {
		user := User{
			Name:  "",                  // First error
			Email: "invalid-email",     // Second error
			Role:  UserRole("INVALID"), // Third error
		}

		err := user.Validate()
		assert.Error(t, err)
		// Should fail on first validation error (name)
		assert.Contains(t, err.Error(), "name is required")
	})
}

// TestValidUserRoles tests the helper function
func TestValidUserRoles(t *testing.T) {
	roles := ValidUserRoles()

	assert.Len(t, roles, 3)
	assert.Contains(t, roles, UserRoleEmployee)
	assert.Contains(t, roles, UserRoleAdmin)
	assert.Contains(t, roles, UserRoleManager)
}
