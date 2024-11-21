package postgres

import (
	"context"
	"database/sql"
	"fmt"

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
	paramCount := 1

	// Add filters
	if filter.Location != nil {
		query += fmt.Sprintf(" AND location = $%d", paramCount)
		countQuery += fmt.Sprintf(" AND location = $%d", paramCount)
		args = append(args, *filter.Location)
		paramCount++
	}
	if filter.JobType != nil {
		query += fmt.Sprintf(" AND job_type = $%d", paramCount)
		countQuery += fmt.Sprintf(" AND job_type = $%d", paramCount)
		args = append(args, *filter.JobType)
		paramCount++
	}
	if filter.ExperienceLevel != nil {
		query += fmt.Sprintf(" AND experience_level = $%d", paramCount)
		countQuery += fmt.Sprintf(" AND experience_level = $%d", paramCount)
		args = append(args, *filter.ExperienceLevel)
		paramCount++
	}
	if len(filter.Skills) > 0 {
		query += fmt.Sprintf(" AND skills && $%d", paramCount)
		countQuery += fmt.Sprintf(" AND skills && $%d", paramCount)
		args = append(args, pq.Array(filter.Skills))
		paramCount++
	}

	// Add pagination
	limit := filter.PageSize
	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, limit, offset)

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
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

func (r *JobRepository) ChangeJobStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
        UPDATE jobs 
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *JobRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Delete associated applications first (assuming you have a foreign key constraint)
	appQuery := `DELETE FROM applications WHERE job_id = $1`
	_, err := r.db.ExecContext(ctx, appQuery, id)
	if err != nil {
		return err
	}

	// Then delete the job
	jobQuery := `DELETE FROM jobs WHERE id = $1`
	_, err = r.db.ExecContext(ctx, jobQuery, id)
	return err
}
