package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type CreateToolRequest struct {
	Name   string            `json:"name" binding:"required"`
	Status domain.ToolStatus `json:"status"`
}

type UpdateToolRequest struct {
	Name   string            `json:"name" binding:"required"`
	Status domain.ToolStatus `json:"status"`
}

// CreateTool godoc
// @Summary Create a new tool
// @Description Create a new tool with name and status
// @Tags tools
// @Accept json
// @Produce json
// @Param tool body CreateToolRequest true "Tool data"
// @Success 201 {object} domain.Tool
// @Failure 400 {object} map[string]string
// @Router /tools [post]
func (s *Server) createTool(c *gin.Context) {
	var req CreateToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status == "" {
		req.Status = domain.ToolStatusInOffice
	}

	tool, err := s.toolService.CreateTool(req.Name, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tool)
}

// ListTools godoc
// @Summary List all tools
// @Description Get a list of tools with pagination and optional status filtering
// @Tags tools
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string][]domain.Tool
// @Failure 400 {object} map[string]string
// @Router /tools [get]
func (s *Server) listTools(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

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

	tools, err := s.toolService.ListTools(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tools"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tools": tools})
}

// GetTool godoc
// @Summary Get a tool by ID
// @Description Get a specific tool by its ID
// @Tags tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Success 200 {object} domain.Tool
// @Failure 404 {object} map[string]string
// @Router /tools/{id} [get]
func (s *Server) getTool(c *gin.Context) {
	id := c.Param("id")

	tool, err := s.toolService.GetTool(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tool not found"})
		return
	}

	c.JSON(http.StatusOK, tool)
}

// UpdateTool godoc
// @Summary Update a tool
// @Description Update a tool's name and status
// @Tags tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Param tool body UpdateToolRequest true "Updated tool data"
// @Success 200 {object} domain.Tool
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tools/{id} [put]
func (s *Server) updateTool(c *gin.Context) {
	id := c.Param("id")

	var req UpdateToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tool, err := s.toolService.UpdateTool(id, req.Name, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tool"})
		return
	}

	c.JSON(http.StatusOK, tool)
}

// DeleteTool godoc
// @Summary Delete a tool
// @Description Delete a tool from the system
// @Tags tools
// @Accept json
// @Produce json
// @Param id path string true "Tool ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tools/{id} [delete]
func (s *Server) deleteTool(c *gin.Context) {
	id := c.Param("id")

	err := s.toolService.DeleteTool(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tool"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
