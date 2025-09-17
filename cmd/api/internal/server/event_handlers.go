package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

// Request structs for tool actions
type CheckoutToolRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Notes  string `json:"notes"`
}

type CheckinToolRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Notes  string `json:"notes"`
}

type MaintenanceRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Notes  string `json:"notes"`
}

type MarkLostRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Notes  string `json:"notes"`
}

// CheckoutTool godoc
// @Summary Check out a tool to a user
// @Description Check out a tool to a specific user with optional notes
// @Tags tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Param checkout body CheckoutToolRequest true "Checkout data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tools/{id}/checkout [post]
func (s *Server) checkoutTool(c *gin.Context) {
	toolID := c.Param("id")

	var req CheckoutToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondDomainError(c, err)
		return
	}

	// Get the tool first to check if it exists and is available
	tool, err := s.toolService.GetTool(toolID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Check if tool is available for checkout
	if tool.Status != domain.ToolStatusInOffice {
		respondDomainError(c, domain.ErrValidation)
		return
	}

	// Update tool status to checked out
	updatedTool, err := s.toolService.UpdateTool(toolID, tool.Name, domain.ToolStatusCheckedOut)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Log the checkout event
	err = s.eventService.LogToolCheckedOut(toolID, req.UserID, req.UserID, req.Notes)
	if err != nil {
		// Log error but don't fail the request
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tool checked out successfully",
		"tool":    updatedTool,
	})
}

// CheckinTool godoc
// @Summary Check in a tool from a user
// @Description Check in a tool that was previously checked out
// @Tags tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Param checkin body CheckinToolRequest true "Checkin data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tools/{id}/checkin [post]
func (s *Server) checkinTool(c *gin.Context) {
	toolID := c.Param("id")

	var req CheckinToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondDomainError(c, err)
		return
	}

	// Get the tool first to check if it exists and is checked out
	tool, err := s.toolService.GetTool(toolID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Check if tool is checked out
	if tool.Status != domain.ToolStatusCheckedOut {
		respondDomainError(c, domain.ErrValidation)
		return
	}

	// Update tool status to in office
	updatedTool, err := s.toolService.UpdateTool(toolID, tool.Name, domain.ToolStatusInOffice)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Log the checkin event
	err = s.eventService.LogToolCheckedIn(toolID, req.UserID, req.UserID, req.Notes)
	if err != nil {
		// Log error but don't fail the request
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tool checked in successfully",
		"tool":    updatedTool,
	})
}

// SendToMaintenance godoc
// @Summary Send a tool to maintenance
// @Description Mark a tool as being in maintenance
// @Tags tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Param maintenance body MaintenanceRequest true "Maintenance data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tools/{id}/maintenance [post]
func (s *Server) sendToMaintenance(c *gin.Context) {
	toolID := c.Param("id")

	var req MaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondDomainError(c, err)
		return
	}

	// Get the tool first
	tool, err := s.toolService.GetTool(toolID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Update tool status to maintenance
	updatedTool, err := s.toolService.UpdateTool(toolID, tool.Name, domain.ToolStatusMaintenance)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Log the maintenance event
	err = s.eventService.LogToolMaintenance(toolID, req.UserID, req.Notes)
	if err != nil {
		// Log error but don't fail the request
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tool sent to maintenance",
		"tool":    updatedTool,
	})
}

// MarkAsLost godoc
// @Summary Mark a tool as lost
// @Description Mark a tool as lost or missing
// @Tags tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Param lost body MarkLostRequest true "Lost data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tools/{id}/lost [post]
func (s *Server) markAsLost(c *gin.Context) {
	toolID := c.Param("id")

	var req MarkLostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondDomainError(c, err)
		return
	}

	// Get the tool first
	tool, err := s.toolService.GetTool(toolID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Update tool status to lost
	updatedTool, err := s.toolService.UpdateTool(toolID, tool.Name, domain.ToolStatusLost)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Log the lost event
	err = s.eventService.LogToolLost(toolID, req.UserID, req.Notes)
	if err != nil {
		// Log error but don't fail the request
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tool marked as lost",
		"tool":    updatedTool,
	})
}

// ListEvents godoc
// @Summary List all events
// @Description Get a list of events with pagination and optional filtering
// @Tags events
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Param type query string false "Filter by event type"
// @Param tool_id query string false "Filter by tool ID"
// @Param user_id query string false "Filter by user ID"
// @Success 200 {object} map[string][]domain.Event
// @Failure 400 {object} map[string]string
// @Router /events [get]
func (s *Server) listEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	eventType := c.Query("type")
	toolID := c.Query("tool_id")
	userID := c.Query("user_id")

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

	// Validate event type if provided
	if eventType != "" && !domain.EventType(eventType).IsValid() {
		respondDomainError(c, domain.ErrValidation)
		return
	}

	var eventTypePtr *string
	var toolIDPtr *string
	var userIDPtr *string

	if eventType != "" {
		eventTypePtr = &eventType
	}
	if toolID != "" {
		toolIDPtr = &toolID
	}
	if userID != "" {
		userIDPtr = &userID
	}

	events, err := s.eventService.ListEvents(limit, offset, eventTypePtr, toolIDPtr, userIDPtr)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetEvent godoc
// @Summary Get an event by ID
// @Description Get a specific event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} domain.Event
// @Failure 404 {object} map[string]string
// @Router /events/{id} [get]
func (s *Server) getEvent(c *gin.Context) {
	id := c.Param("id")

	event, err := s.eventService.GetEvent(id)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetToolHistory godoc
// @Summary Get tool history
// @Description Get the complete event history for a specific tool
// @Tags tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Success 200 {object} map[string][]domain.Event
// @Failure 500 {object} map[string]string
// @Router /tools/{id}/history [get]
func (s *Server) getToolHistory(c *gin.Context) {
	toolID := c.Param("id")

	events, err := s.eventService.GetToolHistory(toolID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetUserActivity godoc
// @Summary Get user activity
// @Description Get the complete activity history for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string][]domain.Event
// @Failure 500 {object} map[string]string
// @Router /users/{id}/activity [get]
func (s *Server) getUserActivity(c *gin.Context) {
	userID := c.Param("id")

	events, err := s.eventService.GetUserActivity(userID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetUserTools godoc
// @Summary Get tools assigned to user
// @Description Get list of tools currently checked out by a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/tools [get]
func (s *Server) getUserTools(c *gin.Context) {
	userID := c.Param("id")

	// This would require a new method in tool service to get tools by user
	// For now, we'll use events to find checked out tools by this user
	events, err := s.eventService.GetUserActivity(userID)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	// Filter for checked out tools (this is a simplified approach)
	var checkedOutTools []string
	for _, event := range events {
		if event.Type == domain.EventTypeToolCheckedOut && event.ToolID != nil {
			checkedOutTools = append(checkedOutTools, *event.ToolID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"checked_out_tool_ids": checkedOutTools})
}
