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