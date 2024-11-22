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

type SocialLinks struct {
	LinkedIn *string `json:"linkedin,omitempty"`
	Twitter  *string `json:"twitter,omitempty"`
	GitHub   *string `json:"github,omitempty"`
}

type ContactInfo struct {
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty"`
	Address *string `json:"address,omitempty"`
}

type EmploymentHistory struct {
	ID          uuid.UUID  `json:"id"`
	Company     string     `json:"company"`
	Title       string     `json:"title"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description *string    `json:"description,omitempty"`
}

type EducationHistory struct {
	ID          uuid.UUID  `json:"id"`
	Institution string     `json:"institution"`
	Degree      string     `json:"degree"`
	Field       string     `json:"field"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type Certification struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Authority    string     `json:"authority"`
	IssueDate    time.Time  `json:"issue_date"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
	CredentialID *string    `json:"credential_id,omitempty"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	FullName     string    `json:"full_name"`
	CompanyName  *string   `json:"company_name,omitempty"`
	ResumeURL    *string   `json:"resume_url,omitempty"`
	Verified     bool      `json:"verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Optional profile details
	Skills            []string            `json:"skills,omitempty"`
	Experience        *string             `json:"experience,omitempty"`
	Education         *string             `json:"education,omitempty"`
	Bio               *string             `json:"bio,omitempty"`
	ProfilePictureURL *string             `json:"profile_picture_url,omitempty"`
	Location          *string             `json:"location,omitempty"`
	SocialLinks       *SocialLinks        `json:"social_links,omitempty"`
	ContactInfo       *ContactInfo        `json:"contact_info,omitempty"`
	EmploymentHistory []EmploymentHistory `json:"employment_history,omitempty"`
	EducationHistory  []EducationHistory  `json:"education_history,omitempty"`
	Certifications    []Certification     `json:"certifications,omitempty"`
}

type UserAnalytics struct {
	UserID              uuid.UUID `json:"user_id"`
	ProfileViews        int       `json:"profile_views"`
	ProfileCompleteness float64   `json:"profile_completeness"`
	LastActive          time.Time `json:"last_active"`
}

type RegisterRequest struct {
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	Role        UserRole `json:"role" binding:"required"`
	FullName    string   `json:"full_name" binding:"required"`
	CompanyName *string  `json:"company_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileDetailsRequest struct {
	Skills            []string     `json:"skills,omitempty"`
	Experience        *string      `json:"experience,omitempty"`
	Education         *string      `json:"education,omitempty"`
	Bio               *string      `json:"bio,omitempty"`
	ProfilePictureURL *string      `json:"profile_picture_url,omitempty"`
	Location          *string      `json:"location,omitempty"`
	SocialLinks       *SocialLinks `json:"social_links,omitempty"`
	ContactInfo       *ContactInfo `json:"contact_info,omitempty"`
}
