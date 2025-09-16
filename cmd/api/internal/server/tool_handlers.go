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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tool"})
		return
	}

	c.JSON(http.StatusCreated, tool)
}

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

func (s *Server) getTool(c *gin.Context) {
	id := c.Param("id")

	tool, err := s.toolService.GetTool(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tool not found"})
		return
	}

	c.JSON(http.StatusOK, tool)
}

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

func (s *Server) deleteTool(c *gin.Context) {
	id := c.Param("id")

	err := s.toolService.DeleteTool(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tool"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
