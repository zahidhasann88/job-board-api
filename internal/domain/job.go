package domain

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	CompanyID       uuid.UUID `json:"company_id"`
	Location        string    `json:"location"`
	SalaryRange     *string   `json:"salary_range,omitempty"`
	JobType         string    `json:"job_type"`
	ExperienceLevel string    `json:"experience_level"`
	Skills          []string  `json:"skills"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type JobFilter struct {
	Location        *string
	JobType         *string
	ExperienceLevel *string
	Skills          []string
	CompanyID       *uuid.UUID
	Status          *string
	Page            int
	PageSize        int
}
type CreateJobRequest struct {
	Title           string   `json:"title" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	Location        string   `json:"location" binding:"required"`
	SalaryRange     *string  `json:"salary_range"`
	JobType         string   `json:"job_type" binding:"required"`
	ExperienceLevel string   `json:"experience_level" binding:"required"`
	Skills          []string `json:"skills" binding:"required"`
}

type UpdateJobRequest struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Location        string   `json:"location"`
	SalaryRange     *string  `json:"salary_range"`
	JobType         string   `json:"job_type"`
	ExperienceLevel string   `json:"experience_level"`
	Skills          []string `json:"skills"`
	Status          string   `json:"status"`
}

type JobAnalytics struct {
    ActiveJobs             int            `json:"active_jobs"`
    TotalApplications      int            `json:"total_applications"`
    ApplicationsByStatus   map[string]int `json:"applications_by_status"`
}

type JobApplicationInsights struct {
    JobID                  uuid.UUID `json:"job_id"`
    TotalApplications      int       `json:"total_applications"`
    PendingApplications    int       `json:"pending_applications"`
    ReviewedApplications   int       `json:"reviewed_applications"`
    InterviewedApplications int      `json:"interviewed_applications"`
    AcceptedApplications   int       `json:"accepted_applications"`
    RejectedApplications   int       `json:"rejected_applications"`
    AverageExperience     float64   `json:"average_experience"`
}

type RecommendedCandidate struct {
    UserID              uuid.UUID `json:"user_id"`
    Name                string    `json:"name"`
    Email               string    `json:"email"`
    YearsOfExperience   int       `json:"years_of_experience"`
    Skills              []string  `json:"skills"`
    SkillMatchPercentage float64  `json:"skill_match_percentage"`
}