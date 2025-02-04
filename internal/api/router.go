package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
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

	// Rate Limiter Setup
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  int64(s.config.RateLimitRequestsPerMinute),
	}
	store := memory.NewStore()
	limiterInstance := limiter.New(store, rate)
	rateLimiterMiddleware := mgin.NewMiddleware(limiterInstance)

	// Apply rate limiter to all routes
	s.router.Use(rateLimiterMiddleware)

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
			recruiter.PUT("/jobs/:id", s.jobHandler.Update)
			recruiter.PATCH("/jobs/:id/status", s.jobHandler.ChangeStatus)
			recruiter.DELETE("/jobs/:id", s.jobHandler.CompleteDelete)
			recruiter.GET("/jobs/analytics", s.jobHandler.GetJobAnalytics)
			recruiter.POST("/jobs/bulk", s.jobHandler.BulkCreateJobs)
			recruiter.GET("/jobs/:id/application-insights", s.jobHandler.GetJobApplicationInsights)
			recruiter.GET("/jobs/:id/recommended-candidates", s.jobHandler.GetRecommendedCandidates)
		}

		// Job seeker routes
		jobSeeker := auth.Group("")
		jobSeeker.Use(middleware.RequireRole(s.config.JobSeekerRole))
		{
			jobSeeker.POST("/applications", s.applicationHandler.Create)
			jobSeeker.GET("/applications", s.applicationHandler.List)
		}

		auth.PUT("/users/profile", s.userHandler.UpdateProfileDetails)
		auth.PUT("/users/employment-history", s.userHandler.UpdateEmploymentHistory)
	}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}
