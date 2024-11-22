package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

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

func (r *UserRepository) CalculateProfileCompleteness(ctx context.Context, userID uuid.UUID) (float64, error) {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	completeness := 20.0 // Base 20% for basic profile

	if user.Skills != nil && len(user.Skills) > 0 {
		completeness += 15
	}
	if user.Experience != nil && *user.Experience != "" {
		completeness += 15
	}
	if user.Education != nil && *user.Education != "" {
		completeness += 15
	}
	if user.Bio != nil && *user.Bio != "" {
		completeness += 10
	}
	if user.ProfilePictureURL != nil {
		completeness += 10
	}
	if len(user.EmploymentHistory) > 0 {
		completeness += 15
	}
	if len(user.EducationHistory) > 0 {
		completeness += 15
	}

	// Update user analytics
	updateQuery := `
        INSERT INTO user_analytics (user_id, profile_completeness)
        VALUES ($1, $2)
        ON CONFLICT (user_id) DO UPDATE 
        SET profile_completeness = $2`

	_, err = r.db.ExecContext(ctx, updateQuery, userID, completeness)
	return completeness, err
}

func (r *UserRepository) UpdateProfileDetails(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	var setClauses []string
	var args []interface{}
	argCounter := 1

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argCounter))
		args = append(args, value)
		argCounter++
	}

	// Always update updated_at
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argCounter))
	args = append(args, time.Now())
	argCounter++

	query := fmt.Sprintf(`
        UPDATE users
        SET %s
        WHERE id = $%d
    `, strings.Join(setClauses, ", "), argCounter)

	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepository) GetUserAnalytics(ctx context.Context, userID uuid.UUID) (*domain.UserAnalytics, error) {
	analytics := &domain.UserAnalytics{
		UserID: userID,
	}
	query := `
		SELECT 
			profile_views, 
			profile_completeness, 
			last_active
		FROM user_analytics
		WHERE user_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&analytics.ProfileViews,
		&analytics.ProfileCompleteness,
		&analytics.LastActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

func (r *UserRepository) IncrementProfileView(ctx context.Context, userID uuid.UUID) error {
	query := `
		INSERT INTO user_analytics (user_id, profile_views, last_active)
		VALUES ($1, 1, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id) DO UPDATE 
		SET profile_views = user_analytics.profile_views + 1,
			last_active = CURRENT_TIMESTAMP
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *UserRepository) UpdateCertifications(ctx context.Context, userID uuid.UUID, certifications []domain.Certification) error {
	// Begin a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, delete existing certifications for the user
	_, err = tx.ExecContext(ctx, "DELETE FROM user_certifications WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Insert new certifications
	for _, cert := range certifications {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO user_certifications (
                id, user_id, name, authority, issue_date, expiry_date, credential_id
            ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        `,
			uuid.New(),
			userID,
			cert.Name,
			cert.Authority,
			cert.IssueDate,
			cert.ExpiryDate,
			cert.CredentialID,
		)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (r *UserRepository) UpdateEducationHistory(ctx context.Context, userID uuid.UUID, history []domain.EducationHistory) error {
	// Begin a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, delete existing education history for the user
	_, err = tx.ExecContext(ctx, "DELETE FROM user_education_history WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Insert new education history entries
	for _, edu := range history {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO user_education_history (
                id, user_id, institution, degree, field, start_date, end_date
            ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        `,
			uuid.New(),
			userID,
			edu.Institution,
			edu.Degree,
			edu.Field,
			edu.StartDate,
			edu.EndDate,
		)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (r *UserRepository) UpdateEmploymentHistory(ctx context.Context, userID uuid.UUID, history []domain.EmploymentHistory) error {
	// Begin a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, delete existing employment history for the user
	_, err = tx.ExecContext(ctx, "DELETE FROM user_employment_history WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Insert new employment history entries
	for _, emp := range history {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO user_employment_history (
                id, user_id, company, title, start_date, end_date, description
            ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        `,
			uuid.New(),
			userID,
			emp.Company,
			emp.Title,
			emp.StartDate,
			emp.EndDate,
			emp.Description,
		)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}
