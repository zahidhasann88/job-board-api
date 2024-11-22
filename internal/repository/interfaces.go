package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job *domain.Job) error
	Update(ctx context.Context, job *domain.Job) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error)
	List(ctx context.Context, filter domain.JobFilter) ([]domain.Job, int, error)
	ChangeJobStatus(ctx context.Context, id uuid.UUID, status string) error
	GetJobAnalytics(ctx context.Context, companyID uuid.UUID) (*domain.JobAnalytics, error)
	GetApplicationInsights(ctx context.Context, jobID uuid.UUID) (*domain.JobApplicationInsights, error)
	BulkCreate(ctx context.Context, jobs []domain.Job) error
	GetRecommendedCandidates(ctx context.Context, jobID uuid.UUID) ([]domain.RecommendedCandidate, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Profile-specific methods
	UpdateProfileDetails(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
	CalculateProfileCompleteness(ctx context.Context, userID uuid.UUID) (float64, error)

	// Detailed history methods
	UpdateEmploymentHistory(ctx context.Context, userID uuid.UUID, history []domain.EmploymentHistory) error
	UpdateEducationHistory(ctx context.Context, userID uuid.UUID, history []domain.EducationHistory) error
	UpdateCertifications(ctx context.Context, userID uuid.UUID, certifications []domain.Certification) error

	// Analytics
	GetUserAnalytics(ctx context.Context, userID uuid.UUID) (*domain.UserAnalytics, error)
	IncrementProfileView(ctx context.Context, userID uuid.UUID) error
}

type ApplicationRepository interface {
	Create(ctx context.Context, application *domain.Application) error
	Update(ctx context.Context, application *domain.Application) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error)
	List(ctx context.Context, filter domain.ApplicationFilter) ([]domain.Application, int, error)

	// GetApplicationsByJob(ctx context.Context, jobID uuid.UUID) ([]domain.Application, error)
	// GetApplicationsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Application, error)
}
