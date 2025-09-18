// @title Tool Tracker API
// @version 1.0
// @description A REST API for tracking company tools and equipment
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/database"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/repo"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/server"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/service"
	_ "github.com/wassaaa/tool-tracker/docs" // This will be generated
)

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://app:app@localhost:5432/tooltracker?sslmode=disable"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database", err)
	}
	log.Println("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Migrations completed")

	toolRepo := repo.NewPostgresToolRepo(db)
	userRepo := repo.NewPostgresUserRepo(db)
	eventRepo := repo.NewPostgresEventRepo(db)

	eventService := service.NewEventService(eventRepo)
	toolService := service.NewToolService(toolRepo).WithEventLogger(eventService)
	userService := service.NewUserService(userRepo).WithEventLogger(eventService)

	srv := server.NewServer(toolService, userService, eventService)

	r := srv.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
