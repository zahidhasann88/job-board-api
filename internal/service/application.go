package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
	"github.com/zahidhasann88/job-board-api/internal/repository"
)

type ApplicationService struct {
	applicationRepo repository.ApplicationRepository
}

func NewApplicationService(applicationRepo repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{applicationRepo: applicationRepo}
}

func (s *ApplicationService) Create(ctx context.Context, application *domain.Application) error {
	application.ID = uuid.New()
	return s.applicationRepo.Create(ctx, application)
}

func (s *ApplicationService) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Application, error) {
	return s.applicationRepo.ListByUser(ctx, userID)
}
