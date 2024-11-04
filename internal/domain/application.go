package domain

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID `json:"id"`
	JobID       uuid.UUID `json:"job_id"`
	ApplicantID uuid.UUID `json:"applicant_id"`
	CoverLetter string    `json:"cover_letter"`
	ResumeURL   string    `json:"resume_url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
