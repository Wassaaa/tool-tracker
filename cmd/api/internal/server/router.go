package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/service"
)

type Server struct {
	toolService *service.ToolService
}

func NewServer(toolService *service.ToolService) *Server {
	return &Server{toolService: toolService}
}

func (s *Server) SetupRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := r.Group("/api")
	api.GET("/tools", s.listTools)
	api.POST("/tools", s.createTool)
	return r
}

type CreateToolRequest struct {
	Name string `json:"name" binding:"required"`
}

func (s *Server) createTool(c *gin.Context) {
	var req CreateToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tool, err := s.toolService.CreateTool(req.Name)
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
