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
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// UpdateSkills(ctx context.Context, userID uuid.UUID, skills []string) error
	// UpdateExperience(ctx context.Context, userID uuid.UUID, experience string) error
	// UpdateEducation(ctx context.Context, userID uuid.UUID, education string) error
	// UpdateBio(ctx context.Context, userID uuid.UUID, bio string) error
	// UpdateProfilePicture(ctx context.Context, userID uuid.UUID, url string) error
	// UpdateResume(ctx context.Context, userID uuid.UUID, url string) error
	// UpdateLocation(ctx context.Context, userID uuid.UUID, location string) error
	// UpdateSocialLinks(ctx context.Context, userID uuid.UUID, links domain.SocialLinks) error
	// UpdateContactInfo(ctx context.Context, userID uuid.UUID, contactInfo domain.ContactInfo) error
	// UpdateEmploymentHistory(ctx context.Context, userID uuid.UUID, employmentHistory []domain.EmploymentHistory) error
	// UpdateEducationHistory(ctx context.Context, userID uuid.UUID, educationHistory []domain.EducationHistory) error
	// UpdateCertifications(ctx context.Context, userID uuid.UUID, certifications []domain.Certification) error
	// GetUserAnalytics(ctx context.Context, userID uuid.UUID) (*domain.UserAnalytics, error)
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
