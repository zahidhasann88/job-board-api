package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zahidhasann88/job-board-api/internal/domain"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job *domain.Job) error {
	query := `
        INSERT INTO jobs (
            id, title, description, company_id, location, salary_range,
            job_type, experience_level, skills, status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		job.ID,
		job.Title,
		job.Description,
		job.CompanyID,
		job.Location,
		job.SalaryRange,
		job.JobType,
		job.ExperienceLevel,
		pq.Array(job.Skills),
		job.Status,
	).Scan(&job.CreatedAt, &job.UpdatedAt)
}

func (r *JobRepository) Update(ctx context.Context, job *domain.Job) error {
	query := `
        UPDATE jobs 
        SET title = $1, description = $2, location = $3, salary_range = $4,
            job_type = $5, experience_level = $6, skills = $7, status = $8,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $9
        RETURNING updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		job.Title,
		job.Description,
		job.Location,
		job.SalaryRange,
		job.JobType,
		job.ExperienceLevel,
		pq.Array(job.Skills),
		job.Status,
		job.ID,
	).Scan(&job.UpdatedAt)
}

func (r *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	job := &domain.Job{}
	query := `
        SELECT id, title, description, company_id, location, salary_range,
               job_type, experience_level, skills, status, created_at, updated_at
        FROM jobs
        WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.Title,
		&job.Description,
		&job.CompanyID,
		&job.Location,
		&job.SalaryRange,
		&job.JobType,
		&job.ExperienceLevel,
		pq.Array(&job.Skills),
		&job.Status,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r *JobRepository) List(ctx context.Context, filter domain.JobFilter) ([]domain.Job, int, error) {
	// Base query
	query := `
        SELECT id, title, description, company_id, location, salary_range,
               job_type, experience_level, skills, status, created_at, updated_at
        FROM jobs
        WHERE 1=1`

	countQuery := "SELECT COUNT(*) FROM jobs WHERE 1=1"
	args := []interface{}{}
	argPosition := 1

	// Add filters
	if filter.Location != nil {
		query += ` AND location = $` + string(argPosition)
		countQuery += ` AND location = $` + string(argPosition)
		args = append(args, *filter.Location)
		argPosition++
	}
	if filter.JobType != nil {
		query += ` AND job_type = $` + string(argPosition)
		countQuery += ` AND job_type = $` + string(argPosition)
		args = append(args, *filter.JobType)
		argPosition++
	}
	if filter.ExperienceLevel != nil {
		query += ` AND experience_level = $` + string(argPosition)
		countQuery += ` AND experience_level = $` + string(argPosition)
		args = append(args, *filter.ExperienceLevel)
		argPosition++
	}
	if len(filter.Skills) > 0 {
		query += ` AND skills && $` + string(argPosition)
		countQuery += ` AND skills && $` + string(argPosition)
		args = append(args, pq.Array(filter.Skills))
		argPosition++
	}

	// Add pagination
	limit := filter.PageSize
	offset := (filter.Page - 1) * filter.PageSize
	query += ` LIMIT $` + string(argPosition) + ` OFFSET $` + string(argPosition+1)
	args = append(args, limit, offset)

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args[:argPosition-1]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Execute main query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []domain.Job
	for rows.Next() {
		var job domain.Job
		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Description,
			&job.CompanyID,
			&job.Location,
			&job.SalaryRange,
			&job.JobType,
			&job.ExperienceLevel,
			pq.Array(&job.Skills),
			&job.Status,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}

	return jobs, total, nil
}
