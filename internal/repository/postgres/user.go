// internal/repository/postgres/user.go
package postgres

import (
	"context"
	"database/sql"

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
        RETURNING created_at, updated_at`

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

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
        SELECT id, email, password_hash, role, full_name, company_name, 
               resume_url, created_at, updated_at
        FROM users
        WHERE email = $1`

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
