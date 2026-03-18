package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/pkg/database"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *database.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, name, password_hash, auth_provider)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Email, user.Name, user.PasswordHash, user.AuthProvider,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, avatar_url, auth_provider, created_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.AvatarURL, &user.AuthProvider, &user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, avatar_url, auth_provider, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.AvatarURL, &user.AuthProvider, &user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = $2, avatar_url = $3
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.AvatarURL)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// CreateOrUpdateProfile creates or updates user profile
func (r *UserRepository) CreateOrUpdateProfile(ctx context.Context, profile *models.UserProfile) error {
	query := `
		INSERT INTO user_profiles (
			user_id, current_job_id, city_id, years_exp,
			education, target_job_id, skills
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id)
		DO UPDATE SET
			current_job_id = EXCLUDED.current_job_id,
			city_id = EXCLUDED.city_id,
			years_exp = EXCLUDED.years_exp,
			education = EXCLUDED.education,
			target_job_id = EXCLUDED.target_job_id,
			skills = EXCLUDED.skills,
			updated_at = NOW()
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		profile.UserID, profile.CurrentJobID, profile.CityID, profile.YearsExp,
		profile.Education, profile.TargetJobID, profile.Skills,
	).Scan(&profile.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create or update profile: %w", err)
	}

	return nil
}

// GetUserProfile retrieves user profile
func (r *UserRepository) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	query := `
		SELECT user_id, current_job_id, city_id, years_exp,
		       education, target_job_id, skills, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	var profile models.UserProfile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.UserID, &profile.CurrentJobID, &profile.CityID, &profile.YearsExp,
		&profile.Education, &profile.TargetJobID, &profile.Skills, &profile.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Profile not found
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &profile, nil
}

// GetUserWithProfile retrieves user with profile
func (r *UserRepository) GetUserWithProfile(ctx context.Context, userID uuid.UUID) (*models.UserWithProfile, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}

	profile, err := r.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.UserWithProfile{
		User:    *user,
		Profile: profile,
	}, nil
}

// SaveAnalysis saves an AI analysis to cache
func (r *UserRepository) SaveAnalysis(ctx context.Context, analysis *models.AIAnalysis) error {
	query := `
		INSERT INTO ai_analyses (
			user_id, profession_slug, location, request_hash,
			result, model_used, tokens_used
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		analysis.UserID, analysis.ProfessionSlug, analysis.Location,
		analysis.RequestHash, analysis.Result, analysis.ModelUsed, analysis.TokensUsed,
	).Scan(&analysis.ID, &analysis.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save analysis: %w", err)
	}

	return nil
}

// GetAnalysisByHash retrieves cached analysis by request hash
func (r *UserRepository) GetAnalysisByHash(ctx context.Context, hash string) (*models.AIAnalysis, error) {
	query := `
		SELECT id, user_id, profession_slug, location, request_hash,
		       result, model_used, tokens_used, created_at
		FROM ai_analyses
		WHERE request_hash = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var analysis models.AIAnalysis
	err := r.db.QueryRow(ctx, query, hash).Scan(
		&analysis.ID, &analysis.UserID, &analysis.ProfessionSlug, &analysis.Location,
		&analysis.RequestHash, &analysis.Result, &analysis.ModelUsed,
		&analysis.TokensUsed, &analysis.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Analysis not found
		}
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	return &analysis, nil
}
