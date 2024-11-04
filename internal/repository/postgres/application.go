package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
)

type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(ctx context.Context, application *domain.Application) error {
	query := `
        INSERT INTO applications (
            id, job_id, applicant_id, cover_letter, resume_url, status
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		application.ID,
		application.JobID,
		application.ApplicantID,
		application.CoverLetter,
		application.ResumeURL,
		application.Status,
	).Scan(&application.CreatedAt, &application.UpdatedAt)
}

func (r *ApplicationRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Application, error) {
	query := `
        SELECT id, job_id, applicant_id, cover_letter, resume_url, status, created_at, updated_at
        FROM applications
        WHERE applicant_id = $1
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []domain.Application
	for rows.Next() {
		var app domain.Application
		if err := rows.Scan(
			&app.ID,
			&app.JobID,
			&app.ApplicantID,
			&app.CoverLetter,
			&app.ResumeURL,
			&app.Status,
			&app.CreatedAt,
			&app.UpdatedAt,
		); err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applications, nil
}
