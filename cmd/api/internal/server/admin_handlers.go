package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Admin stats response
type StatsResponse struct {
	TotalTools    int `json:"total_tools"`
	TotalUsers    int `json:"total_users"`
	TotalEvents   int `json:"total_events"`
	ToolsByStatus struct {
		InOffice    int `json:"in_office"`
		CheckedOut  int `json:"checked_out"`
		Maintenance int `json:"maintenance"`
		Lost        int `json:"lost"`
	} `json:"tools_by_status"`
	UsersByRole struct {
		Employees int `json:"employees"`
		Managers  int `json:"managers"`
		Admins    int `json:"admins"`
	} `json:"users_by_role"`
}

// GetStats godoc
// @Summary Get system statistics
// @Description Get comprehensive statistics about tools, users, and events
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} StatsResponse
// @Failure 500 {object} map[string]string
// @Router /admin/stats [get]
func (s *Server) getStats(c *gin.Context) {
	var stats StatsResponse

	// Get total counts
	toolCount, err := s.toolService.GetToolCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tool count"})
		return
	}
	stats.TotalTools = toolCount

	userCount, err := s.userService.GetUserCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count"})
		return
	}
	stats.TotalUsers = userCount

	eventCount, err := s.eventService.GetEventCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get event count"})
		return
	}
	stats.TotalEvents = eventCount

	// Note: For tools by status and users by role, you would need to add
	// specific methods to your services or repositories to count by these criteria
	// For now, we'll leave them as 0

	c.JSON(http.StatusOK, stats)
}

// GetAuditLog godoc
// @Summary Get audit log
// @Description Get recent audit events for administrative review
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /admin/audit [get]
func (s *Server) getAuditLog(c *gin.Context) {
	// Get recent audit events (last 100)
	events, err := s.eventService.ListEvents(100, 0, nil, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get audit log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audit_log": events,
		"total":     len(events),
	})
}
