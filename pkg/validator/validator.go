// pkg/validator/validator.go
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CustomValidator struct {
	validator *validator.Validate
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var errMsgs []string
	for _, err := range ve {
		errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(errMsgs, "; ")
}

// NewValidator creates a new custom validator
func NewValidator() *CustomValidator {
	v := validator.New()

	// Register custom validation functions
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("uuid", validateUUID)
	v.RegisterValidation("job_type", validateJobType)
	v.RegisterValidation("experience_level", validateExperienceLevel)
	v.RegisterValidation("application_status", validateApplicationStatus)
	v.RegisterValidation("url", validateURL)
	v.RegisterValidation("salary_range", validateSalaryRange)

	return &CustomValidator{
		validator: v,
	}
}

// Validate validates the input struct and returns ValidationErrors
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors

	for _, err := range err.(validator.ValidationErrors) {
		field := toSnakeCase(err.Field())
		message := getErrorMessage(err)

		validationErrors = append(validationErrors, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	return validationErrors
}

// Custom validation functions
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Password must be at least 8 characters long and contain at least:
	// 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return len(password) >= 8 && hasUpper && hasLower && hasNumber && hasSpecial
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Basic phone number validation (can be customized based on your requirements)
	match, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, phone)
	return match
}

func validateUUID(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	_, err := uuid.Parse(id)
	return err == nil
}

func validateJobType(fl validator.FieldLevel) bool {
	validTypes := map[string]bool{
		"full-time":  true,
		"part-time":  true,
		"contract":   true,
		"internship": true,
		"freelance":  true,
		"remote":     true,
	}
	return validTypes[strings.ToLower(fl.Field().String())]
}

func validateExperienceLevel(fl validator.FieldLevel) bool {
	validLevels := map[string]bool{
		"entry":     true,
		"junior":    true,
		"mid":       true,
		"senior":    true,
		"lead":      true,
		"executive": true,
	}
	return validLevels[strings.ToLower(fl.Field().String())]
}

func validateApplicationStatus(fl validator.FieldLevel) bool {
	validStatus := map[string]bool{
		"pending":  true,
		"reviewed": true,
		"accepted": true,
		"rejected": true,
	}
	return validStatus[strings.ToLower(fl.Field().String())]
}

func validateURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	match, _ := regexp.MatchString(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/\S*)?$`, url)
	return match
}

func validateSalaryRange(fl validator.FieldLevel) bool {
	salaryRange := fl.Field().String()
	// Matches patterns like "50000-75000" or "50k-75k" or "50K-75K"
	match, _ := regexp.MatchString(`^\d+k?-\d+k?$`, strings.ToLower(salaryRange))
	return match
}

// Helper functions
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "password":
		return "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"
	case "phone":
		return "Invalid phone number format"
	case "uuid":
		return "Invalid UUID format"
	case "job_type":
		return "Invalid job type. Must be one of: full-time, part-time, contract, internship, freelance, remote"
	case "experience_level":
		return "Invalid experience level. Must be one of: entry, junior, mid, senior, lead, executive"
	case "application_status":
		return "Invalid application status. Must be one of: pending, reviewed, accepted, rejected"
	case "url":
		return "Invalid URL format"
	case "salary_range":
		return "Invalid salary range format. Example: 50000-75000 or 50k-75k"
	default:
		return fmt.Sprintf("Failed validation on %s", err.Tag())
	}
}

// Example usage in request structs
type CreateJobRequest struct {
	Title           string   `json:"title" validate:"required,min=3,max=100"`
	Description     string   `json:"description" validate:"required,min=10"`
	Location        string   `json:"location" validate:"required"`
	SalaryRange     string   `json:"salary_range" validate:"required,salary_range"`
	JobType         string   `json:"job_type" validate:"required,job_type"`
	ExperienceLevel string   `json:"experience_level" validate:"required,experience_level"`
	Skills          []string `json:"skills" validate:"required,min=1,dive,required"`
}

type ApplicationRequest struct {
	JobID       string `json:"job_id" validate:"required,uuid"`
	CoverLetter string `json:"cover_letter" validate:"required,min=50"`
	ResumeURL   string `json:"resume_url" validate:"required,url"`
}

type UpdateProfileRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=2,max=100"`
	Phone       string  `json:"phone" validate:"required,phone"`
	CompanyName *string `json:"company_name,omitempty" validate:"omitempty,min=2,max=100"`
	ResumeURL   *string `json:"resume_url,omitempty" validate:"omitempty,url"`
}
