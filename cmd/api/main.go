package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/zahidhasann88/job-board-api/internal/api"
	"github.com/zahidhasann88/job-board-api/internal/config"
	"github.com/zahidhasann88/job-board-api/internal/repository/postgres"
	"github.com/zahidhasann88/job-board-api/internal/service"
	"github.com/zahidhasann88/job-board-api/pkg/logger"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	l := logger.NewLogger()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	jobRepo := postgres.NewJobRepository(db)
	applicationRepo := postgres.NewApplicationRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	jobService := service.NewJobService(jobRepo)
	applicationService := service.NewApplicationService(applicationRepo)

	// Initialize and start the server
	server := api.NewServer(cfg, l, userService, jobService, applicationService)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
