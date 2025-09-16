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
	toolService := service.NewToolService(toolRepo)
	srv := server.NewServer(toolService)

	r := srv.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
