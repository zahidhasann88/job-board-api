package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleJobSeeker UserRole = "job_seeker"
	RoleRecruiter UserRole = "recruiter"
	RoleAdmin     UserRole = "admin"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	FullName     string    `json:"full_name"`
	CompanyName  *string   `json:"company_name,omitempty"`
	ResumeURL    *string   `json:"resume_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
