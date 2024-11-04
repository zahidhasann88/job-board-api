package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/zahidhasann88/job-board-api/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (
            id, email, password_hash, role, full_name, company_name, resume_url
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at, updated_at
    `

	return r.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.FullName,
		user.CompanyName,
		user.ResumeURL,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	query := `
        SELECT id, email, password_hash, role, full_name, company_name, 
               resume_url, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.FullName,
		&user.CompanyName,
		&user.ResumeURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
        SELECT id, email, password_hash, role, full_name, company_name, 
               resume_url, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.FullName,
		&user.CompanyName,
		&user.ResumeURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users
        SET full_name = $1, company_name = $2, resume_url = $3, updated_at = CURRENT_TIMESTAMP
        WHERE id = $4
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.FullName,
		user.CompanyName,
		user.ResumeURL,
		user.ID,
	)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM users
        WHERE id = $1
    `

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}