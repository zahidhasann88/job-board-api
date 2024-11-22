package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
	"github.com/zahidhasann88/job-board-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User, password string) error {
	// Check if user exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.ID = uuid.New()
	user.PasswordHash = string(hashedPassword)
	user.Verified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) UpdateProfileDetails(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	// Validate updates
	allowedFields := map[string]bool{
		"skills":              true,
		"experience":          true,
		"education":           true,
		"bio":                 true,
		"profile_picture_url": true,
		"location":            true,
		"social_links":        true,
		"contact_info":        true,
	}

	for key := range updates {
		if !allowedFields[key] {
			return fmt.Errorf("invalid field: %s", key)
		}
	}

	// Update profile details
	if err := s.userRepo.UpdateProfileDetails(ctx, userID, updates); err != nil {
		return err
	}

	// Recalculate profile completeness
	_, err := s.userRepo.CalculateProfileCompleteness(ctx, userID)
	return err
}

func (s *UserService) UpdateEmploymentHistory(ctx context.Context, userID uuid.UUID, history []domain.EmploymentHistory) error {
	return s.userRepo.UpdateEmploymentHistory(ctx, userID, history)
}

func (s *UserService) UpdateEducationHistory(ctx context.Context, userID uuid.UUID, history []domain.EducationHistory) error {
	return s.userRepo.UpdateEducationHistory(ctx, userID, history)
}

func (s *UserService) UpdateCertifications(ctx context.Context, userID uuid.UUID, certifications []domain.Certification) error {
	return s.userRepo.UpdateCertifications(ctx, userID, certifications)
}
