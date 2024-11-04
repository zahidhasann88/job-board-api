package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zahidhasann88/job-board-api/internal/api/handler"
	"github.com/zahidhasann88/job-board-api/internal/api/middleware"
	"github.com/zahidhasann88/job-board-api/internal/config"
	"github.com/zahidhasann88/job-board-api/internal/service"
	"go.uber.org/zap"
)

type Server struct {
	config             *config.Config
	logger             *zap.Logger
	router             *gin.Engine
	userHandler        *handler.UserHandler
	jobHandler         *handler.JobHandler
	applicationHandler *handler.ApplicationHandler
}

func NewServer(
	cfg *config.Config,
	logger *zap.Logger,
	userService *service.UserService,
	jobService *service.JobService,
	applicationService *service.ApplicationService,
) *Server {
	server := &Server{
		config:             cfg,
		logger:             logger,
		router:             gin.New(),
		userHandler:        handler.NewUserHandler(userService),
		jobHandler:         handler.NewJobHandler(jobService),
		applicationHandler: handler.NewApplicationHandler(applicationService),
	}
	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	// Middleware
	s.router.Use(middleware.LoggingMiddleware(s.logger))
	s.router.Use(gin.Recovery())

	// Public routes
	s.router.POST("/api/v1/users/register", s.userHandler.Register)
	s.router.POST("/api/v1/users/login", s.userHandler.Login)
	s.router.GET("/api/v1/jobs", s.jobHandler.List)
	s.router.GET("/api/v1/jobs/:id", s.jobHandler.Get)

	// Protected routes
	auth := s.router.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware(s.config.JWTSecret))
	{
		// Recruiter routes
		recruiter := auth.Group("")
		recruiter.Use(middleware.RequireRole(s.config.RecruiterRole))
		{
			recruiter.POST("/jobs", s.jobHandler.Create)
			// recruiter.PUT("/jobs/:id", s.jobHandler.Update)
			// recruiter.PUT("/users/profile", s.userHandler.UpdateProfile)
		}

		// Job seeker routes
		jobSeeker := auth.Group("")
		jobSeeker.Use(middleware.RequireRole(s.config.JobSeekerRole))
		{
			jobSeeker.POST("/applications", s.applicationHandler.Create)
			jobSeeker.GET("/applications", s.applicationHandler.List)
		}

		// Common routes
		auth.PUT("/users/profile", s.userHandler.UpdateProfile)
	}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}
