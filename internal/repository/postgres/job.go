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

func (r *JobRepository) GetJobAnalytics(ctx context.Context, companyID uuid.UUID) (*domain.JobAnalytics, error) {
	analytics := &domain.JobAnalytics{}

	// Get total active jobs
	activeJobsQuery := `
        SELECT COUNT(*) 
        FROM jobs 
        WHERE company_id = $1 AND status = 'active'`
	err := r.db.QueryRowContext(ctx, activeJobsQuery, companyID).Scan(&analytics.ActiveJobs)
	if err != nil {
		return nil, err
	}

	// Get total applications
	applicationsQuery := `
        SELECT COUNT(*) 
        FROM applications a 
        JOIN jobs j ON a.job_id = j.id 
        WHERE j.company_id = $1`
	err = r.db.QueryRowContext(ctx, applicationsQuery, companyID).Scan(&analytics.TotalApplications)
	if err != nil {
		return nil, err
	}

	// Get applications by status
	statusQuery := `
        SELECT a.status, COUNT(*) 
        FROM applications a 
        JOIN jobs j ON a.job_id = j.id 
        WHERE j.company_id = $1 
        GROUP BY a.status`
	rows, err := r.db.QueryContext(ctx, statusQuery, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	analytics.ApplicationsByStatus = make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		analytics.ApplicationsByStatus[status] = count
	}

	return analytics, nil
}

func (r *JobRepository) BulkCreate(ctx context.Context, jobs []domain.Job) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO jobs (
            id, title, description, company_id, location, salary_range,
            job_type, experience_level, skills, status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, job := range jobs {
		_, err = stmt.ExecContext(ctx,
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
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *JobRepository) GetApplicationInsights(ctx context.Context, jobID uuid.UUID) (*domain.JobApplicationInsights, error) {
	insights := &domain.JobApplicationInsights{
		JobID: jobID,
	}

	// Get application statistics
	statsQuery := `
        SELECT 
            COUNT(*) as total_applications,
            COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_applications,
            COUNT(CASE WHEN status = 'reviewed' THEN 1 END) as reviewed_applications,
            COUNT(CASE WHEN status = 'interviewed' THEN 1 END) as interviewed_applications,
            COUNT(CASE WHEN status = 'accepted' THEN 1 END) as accepted_applications,
            COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected_applications
        FROM applications
        WHERE job_id = $1`

	err := r.db.QueryRowContext(ctx, statsQuery, jobID).Scan(
		&insights.TotalApplications,
		&insights.PendingApplications,
		&insights.ReviewedApplications,
		&insights.InterviewedApplications,
		&insights.AcceptedApplications,
		&insights.RejectedApplications,
	)
	if err != nil {
		return nil, err
	}

	// Get average experience of applicants
	expQuery := `
        SELECT AVG(years_of_experience) 
        FROM applications 
        WHERE job_id = $1`
	err = r.db.QueryRowContext(ctx, expQuery, jobID).Scan(&insights.AverageExperience)
	if err != nil {
		return nil, err
	}

	return insights, nil
}

func (r *JobRepository) GetRecommendedCandidates(ctx context.Context, jobID uuid.UUID) ([]domain.RecommendedCandidate, error) {
	query := `
        WITH job_skills AS (
            SELECT skills FROM jobs WHERE id = $1
        )
        SELECT 
            u.id,
            u.name,
            u.email,
            u.years_of_experience,
            array_agg(DISTINCT us.skill) as skills,
            COUNT(DISTINCT us.skill) * 100.0 / (
                SELECT array_length(skills, 1) FROM job_skills
            ) as skill_match_percentage
        FROM 
            users u
            JOIN user_skills us ON u.id = us.user_id
            CROSS JOIN job_skills
            WHERE us.skill = ANY((SELECT skills FROM job_skills))
            AND u.role = 'job_seeker'
            AND NOT EXISTS (
                SELECT 1 FROM applications a 
                WHERE a.user_id = u.id AND a.job_id = $1
            )
        GROUP BY u.id, u.name, u.email, u.years_of_experience
        HAVING COUNT(DISTINCT us.skill) >= (
            SELECT array_length(skills, 1) * 0.5 FROM job_skills
        )
        ORDER BY skill_match_percentage DESC
        LIMIT 10`

	rows, err := r.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []domain.RecommendedCandidate
	for rows.Next() {
		var c domain.RecommendedCandidate
		err := rows.Scan(
			&c.UserID,
			&c.Name,
			&c.Email,
			&c.YearsOfExperience,
			pq.Array(&c.Skills),
			&c.SkillMatchPercentage,
		)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, c)
	}

	return candidates, nil
}
