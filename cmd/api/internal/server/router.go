package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/service"
)

type Server struct {
	toolService  *service.ToolService
	userService  *service.UserService
	eventService *service.EventService
}

func NewServer(
	toolService *service.ToolService,
	userService *service.UserService,
	eventService *service.EventService,
) *Server {
	return &Server{
		toolService:  toolService,
		userService:  userService,
		eventService: eventService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := r.Group("/api")
	{
		// Tools (CRUD)
		tools := api.Group("/tools")
		{
			tools.GET("", s.listTools)
			tools.POST("", s.createTool)
			tools.GET("/:id", s.getTool)
			tools.PUT("/:id", s.updateTool)
			tools.DELETE("/:id", s.deleteTool)

			// Tool Actions (Events)
			tools.POST("/:id/checkout", s.checkoutTool)
			tools.POST("/:id/checkin", s.checkinTool)
			tools.POST("/:id/maintenance", s.sendToMaintenance)
			tools.POST("/:id/lost", s.markAsLost)

			// Tool History
			tools.GET("/:id/history", s.getToolHistory)
		}

		// Users (CRUD)
		users := api.Group("/users")
		{
			users.GET("", s.listUsers)
			users.POST("", s.createUser)
			users.GET("/:id", s.getUser)
			users.PUT("/:id", s.updateUser)
			users.DELETE("/:id", s.deleteUser)

			// User Activity
			users.GET("/:id/activity", s.getUserActivity)
			users.GET("/:id/tools", s.getUserTools)
		}

		// Events/Audit Log
		events := api.Group("/events")
		{
			events.GET("", s.listEvents)
			events.GET("/:id", s.getEvent)
		}

		// Admin routes
		admin := api.Group("/admin")
		{
			admin.GET("/stats", s.getStats)
			admin.GET("/audit", s.getAuditLog)
		}
	}
	return r
}
