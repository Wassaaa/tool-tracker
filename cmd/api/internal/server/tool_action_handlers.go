package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Request payloads for tool actions.
// Logging & state transitions are handled inside the ToolService; handlers just dispatch.
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

	actor := GetActorID(c)
	updatedTool, err := s.toolService.CheckOutTool(toolID, req.UserID, actor, req.Notes)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tool checked out successfully", "tool": updatedTool})
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

	actor := GetActorID(c)
	updatedTool, err := s.toolService.ReturnTool(toolID, actor, req.Notes)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tool checked in successfully", "tool": updatedTool})
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

	actor := GetActorID(c)
	updatedTool, err := s.toolService.SendToMaintenance(toolID, actor, req.Notes)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tool sent to maintenance", "tool": updatedTool})
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

	actor := GetActorID(c)
	updatedTool, err := s.toolService.MarkLost(toolID, actor, req.Notes)
	if err != nil {
		respondDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tool marked as lost", "tool": updatedTool})
}
