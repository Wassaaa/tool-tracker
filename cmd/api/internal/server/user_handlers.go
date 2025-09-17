package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type CreateUserRequest struct {
	Name  string          `json:"name" binding:"required"`
	Email string          `json:"email" binding:"required"`
	Role  domain.UserRole `json:"role"`
}

type UpdateUserRequest struct {
	Name  string          `json:"name" binding:"required"`
	Email string          `json:"email" binding:"required"`
	Role  domain.UserRole `json:"role" binding:"required"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name, email, and role
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (s *Server) createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = domain.UserRoleEmployee
	}

	// Validate role
	if !req.Role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role. Valid values: EMPLOYEE, ADMIN, MANAGER",
		})
		return
	}

	user, err := s.userService.CreateUser(req.Name, req.Email, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of users with pagination and optional role filtering
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param role query string false "Filter by role (EMPLOYEE, ADMIN, MANAGER)"
// @Success 200 {object} map[string][]domain.User
// @Failure 400 {object} map[string]string
// @Router /users [get]
func (s *Server) listUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	roleFilter := c.Query("role") // Optional role filter

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	// Validate role filter if provided
	if roleFilter != "" && !domain.UserRole(roleFilter).IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role filter. Valid values: EMPLOYEE, ADMIN, MANAGER",
		})
		return
	}

	var users []domain.User
	if roleFilter != "" {
		users, err = s.userService.ListUsersByRole(domain.UserRole(roleFilter), limit, offset)
	} else {
		users, err = s.userService.ListUsers(limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (s *Server) getUser(c *gin.Context) {
	id := c.Param("id")

	user, err := s.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's name, email, and role
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserRequest true "Updated user data"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func (s *Server) updateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role
	if !req.Role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role. Valid values: EMPLOYEE, ADMIN, MANAGER",
		})
		return
	}

	user, err := s.userService.UpdateUser(id, req.Name, req.Email, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user from the system
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (s *Server) deleteUser(c *gin.Context) {
	id := c.Param("id")

	err := s.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Helper function to validate user role
func isValidUserRole(role domain.UserRole) bool {
	switch role {
	case domain.UserRoleEmployee, domain.UserRoleAdmin, domain.UserRoleManager:
		return true
	default:
		return false
	}
}
