package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
	"github.com/zahidhasann88/job-board-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	userRepo repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo: userRepo,
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
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, user *domain.User) error {
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	existingUser.FullName = user.FullName
	existingUser.CompanyName = user.CompanyName
	existingUser.ResumeURL = user.ResumeURL

	return s.userRepo.Update(ctx, existingUser)
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