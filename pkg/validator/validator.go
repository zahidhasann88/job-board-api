package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// UserRole represents different user roles in the system
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleRecruiter UserRole = "recruiter"
	RoleApplicant UserRole = "job_seeker"
)

// ValidationContext contains additional context for validation
type ValidationContext struct {
	Role      UserRole
	UserID    string
	CompanyID string
}

type CustomValidator struct {
	validator *validator.Validate
	context   *ValidationContext // Store context in the validator
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
	cv := &CustomValidator{
		validator: v,
	}

	// Register custom validation functions
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("uuid", validateUUID)
	v.RegisterValidation("job_type", validateJobType)
	v.RegisterValidation("experience_level", validateExperienceLevel)
	v.RegisterValidation("application_status", validateApplicationStatus)
	v.RegisterValidation("url", validateURL)
	v.RegisterValidation("salary_range", validateSalaryRange)

	// Register role-based validation functions with closure to access context
	v.RegisterValidation("admin_only", cv.validateAdminOnly)
	v.RegisterValidation("recruiter_only", cv.validateRecruiterOnly)
	v.RegisterValidation("same_user", cv.validateSameUser)
	v.RegisterValidation("same_company", cv.validateSameCompany)

	return cv
}

// Validate validates the input struct without role context
func (cv *CustomValidator) Validate(i interface{}) error {
	cv.context = nil // Clear any existing context
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

// ValidateWithRole validates the input struct with role context
func (cv *CustomValidator) ValidateWithRole(i interface{}, ctx ValidationContext) error {
	// Store context for use in validation functions
	cv.context = &ctx

	err := cv.validator.Struct(i)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors

	for _, err := range err.(validator.ValidationErrors) {
		field := toSnakeCase(err.Field())
		message := getRoleErrorMessage(err, ctx.Role)

		validationErrors = append(validationErrors, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	// Clear context after validation
	cv.context = nil

	return validationErrors
}

// Existing validation functions remain the same
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return len(password) >= 8 && hasUpper && hasLower && hasNumber && hasSpecial
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
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
	match, _ := regexp.MatchString(`^\d+k?-\d+k?$`, strings.ToLower(salaryRange))
	return match
}

// Role-based validation functions with receiver to access context
func (cv *CustomValidator) validateAdminOnly(fl validator.FieldLevel) bool {
	if cv.context == nil {
		return true // If no context, skip role validation
	}
	return cv.context.Role == RoleAdmin
}

func (cv *CustomValidator) validateRecruiterOnly(fl validator.FieldLevel) bool {
	if cv.context == nil {
		return true // If no context, skip role validation
	}
	return cv.context.Role == RoleRecruiter || cv.context.Role == RoleAdmin
}

func (cv *CustomValidator) validateSameUser(fl validator.FieldLevel) bool {
	if cv.context == nil {
		return true // If no context, skip role validation
	}
	return fl.Field().String() == cv.context.UserID || cv.context.Role == RoleAdmin
}

func (cv *CustomValidator) validateSameCompany(fl validator.FieldLevel) bool {
	if cv.context == nil {
		return true // If no context, skip role validation
	}
	return fl.Field().String() == cv.context.CompanyID || cv.context.Role == RoleAdmin
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

func getRoleErrorMessage(err validator.FieldError, role UserRole) string {
	switch err.Tag() {
	case "admin_only":
		return "This field can only be modified by administrators"
	case "recruiter_only":
		return "This field can only be modified by recruiters or administrators"
	case "same_user":
		return "You can only modify your own information"
	case "same_company":
		return "You can only modify information for your own company"
	default:
		return getErrorMessage(err)
	}
}

// Example structs showing usage
type JobPosting struct {
	ID              string `json:"id" validate:"required,uuid"`
	Title           string `json:"title" validate:"required"`
	CompanyID       string `json:"company_id" validate:"required,uuid,same_company"`
	SalaryRange     string `json:"salary_range" validate:"required,salary_range,recruiter_only"`
	FeaturedListing bool   `json:"featured_listing" validate:"admin_only"`
	JobType         string `json:"job_type" validate:"required,job_type"`
	ExperienceLevel string `json:"experience_level" validate:"required,experience_level"`
}

type UserProfile struct {
	ID          string `json:"id" validate:"required,uuid,same_user"`
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Role        string `json:"role" validate:"admin_only"`
	PhoneNumber string `json:"phone_number" validate:"required,phone"`
	Password    string `json:"password" validate:"required,password"`
}

type JobApplication struct {
	ID          string `json:"id" validate:"required,uuid"`
	JobID       string `json:"job_id" validate:"required,uuid"`
	ApplicantID string `json:"applicant_id" validate:"required,uuid,same_user"`
	Status      string `json:"status" validate:"required,application_status,recruiter_only"`
	ResumeURL   string `json:"resume_url" validate:"required,url"`
}
