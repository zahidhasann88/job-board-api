package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
	"github.com/zahidhasann88/job-board-api/internal/repository"
)

type JobService struct {
	jobRepo repository.JobRepository
}

func NewJobService(jobRepo repository.JobRepository) *JobService {
	return &JobService{jobRepo: jobRepo}
}

func (s *JobService) CreateJob(ctx context.Context, job *domain.Job) error {
	job.ID = uuid.New()
	job.Status = "active"
	return s.jobRepo.Create(ctx, job)
}

func (s *JobService) UpdateJob(ctx context.Context, job *domain.Job) error {
	return s.jobRepo.Update(ctx, job)
}

func (s *JobService) GetJob(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	return s.jobRepo.GetByID(ctx, id)
}

func (s *JobService) ListJobs(ctx context.Context, filter domain.JobFilter) ([]domain.Job, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	return s.jobRepo.List(ctx, filter)
}

func (s *JobService) ChangeJobStatus(ctx context.Context, id uuid.UUID, status string) error {
	// Add validation for allowed status values
	allowedStatuses := []string{"active", "inactive", "closed", "draft"}
	if !contains(allowedStatuses, status) {
		return fmt.Errorf("invalid job status: %s", status)
	}
	return s.jobRepo.ChangeJobStatus(ctx, id, status)
}

func (s *JobService) CompleteDeleteJob(ctx context.Context, id uuid.UUID) error {
	return s.jobRepo.Delete(ctx, id)
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
