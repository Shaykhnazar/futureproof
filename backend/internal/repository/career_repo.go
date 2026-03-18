package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/pkg/database"
)

// CareerRepository handles career-related database operations
type CareerRepository struct {
	db *database.DB
}

// NewCareerRepository creates a new career repository
func NewCareerRepository(db *database.DB) *CareerRepository {
	return &CareerRepository{db: db}
}

// GetAllProfessions retrieves all professions
func (r *CareerRepository) GetAllProfessions(ctx context.Context) ([]models.Profession, error) {
	query := `
		SELECT id, slug, title, category, ai_risk_score, avg_salary_usd,
		       description, is_future_job, demand_index, growth_pct, updated_at
		FROM professions
		ORDER BY demand_index DESC, title ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query professions: %w", err)
	}
	defer rows.Close()

	var professions []models.Profession
	for rows.Next() {
		var p models.Profession
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Category, &p.AIRiskScore,
			&p.AvgSalaryUSD, &p.Description, &p.IsFutureJob,
			&p.DemandIndex, &p.GrowthPct, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profession: %w", err)
		}
		professions = append(professions, p)
	}

	return professions, rows.Err()
}

// GetProfessionBySlug retrieves a profession by its slug
func (r *CareerRepository) GetProfessionBySlug(ctx context.Context, slug string) (*models.Profession, error) {
	query := `
		SELECT id, slug, title, category, ai_risk_score, avg_salary_usd,
		       description, is_future_job, demand_index, growth_pct, updated_at
		FROM professions
		WHERE slug = $1
	`

	var p models.Profession
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.Slug, &p.Title, &p.Category, &p.AIRiskScore,
		&p.AvgSalaryUSD, &p.Description, &p.IsFutureJob,
		&p.DemandIndex, &p.GrowthPct, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get profession by slug: %w", err)
	}

	return &p, nil
}

// GetProfessionByID retrieves a profession by its ID
func (r *CareerRepository) GetProfessionByID(ctx context.Context, id uuid.UUID) (*models.Profession, error) {
	query := `
		SELECT id, slug, title, category, ai_risk_score, avg_salary_usd,
		       description, is_future_job, demand_index, growth_pct, updated_at
		FROM professions
		WHERE id = $1
	`

	var p models.Profession
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Slug, &p.Title, &p.Category, &p.AIRiskScore,
		&p.AvgSalaryUSD, &p.Description, &p.IsFutureJob,
		&p.DemandIndex, &p.GrowthPct, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get profession by ID: %w", err)
	}

	return &p, nil
}

// GetCareerTransitions retrieves career pivot recommendations for a profession
func (r *CareerRepository) GetCareerTransitions(ctx context.Context, professionSlug string) ([]models.CareerTransitionWithDetails, error) {
	query := `
		SELECT
			ct.id, ct.match_score, ct.transition_reason, ct.avg_reskill_months,
			fp.id, fp.slug, fp.title, fp.category, fp.ai_risk_score, fp.avg_salary_usd,
			fp.description, fp.is_future_job, fp.demand_index, fp.growth_pct, fp.updated_at,
			tp.id, tp.slug, tp.title, tp.category, tp.ai_risk_score, tp.avg_salary_usd,
			tp.description, tp.is_future_job, tp.demand_index, tp.growth_pct, tp.updated_at
		FROM career_transitions ct
		JOIN professions fp ON ct.from_profession = fp.id
		JOIN professions tp ON ct.to_profession = tp.id
		WHERE fp.slug = $1
		ORDER BY ct.match_score DESC
		LIMIT 10
	`

	rows, err := r.db.Query(ctx, query, professionSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to query career transitions: %w", err)
	}
	defer rows.Close()

	var transitions []models.CareerTransitionWithDetails
	for rows.Next() {
		var t models.CareerTransitionWithDetails
		err := rows.Scan(
			&t.ID, &t.MatchScore, &t.TransitionReason, &t.AvgReskillMonths,
			&t.FromProfession.ID, &t.FromProfession.Slug, &t.FromProfession.Title,
			&t.FromProfession.Category, &t.FromProfession.AIRiskScore, &t.FromProfession.AvgSalaryUSD,
			&t.FromProfession.Description, &t.FromProfession.IsFutureJob, &t.FromProfession.DemandIndex,
			&t.FromProfession.GrowthPct, &t.FromProfession.UpdatedAt,
			&t.ToProfession.ID, &t.ToProfession.Slug, &t.ToProfession.Title,
			&t.ToProfession.Category, &t.ToProfession.AIRiskScore, &t.ToProfession.AvgSalaryUSD,
			&t.ToProfession.Description, &t.ToProfession.IsFutureJob, &t.ToProfession.DemandIndex,
			&t.ToProfession.GrowthPct, &t.ToProfession.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan career transition: %w", err)
		}
		transitions = append(transitions, t)
	}

	return transitions, rows.Err()
}

// GetFutureProfessions retrieves all future/emerging professions
func (r *CareerRepository) GetFutureProfessions(ctx context.Context) ([]models.Profession, error) {
	query := `
		SELECT id, slug, title, category, ai_risk_score, avg_salary_usd,
		       description, is_future_job, demand_index, growth_pct, updated_at
		FROM professions
		WHERE is_future_job = true
		ORDER BY growth_pct DESC, demand_index DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query future professions: %w", err)
	}
	defer rows.Close()

	var professions []models.Profession
	for rows.Next() {
		var p models.Profession
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Category, &p.AIRiskScore,
			&p.AvgSalaryUSD, &p.Description, &p.IsFutureJob,
			&p.DemandIndex, &p.GrowthPct, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profession: %w", err)
		}
		professions = append(professions, p)
	}

	return professions, rows.Err()
}

// SaveCareer saves a profession to user's saved list
func (r *CareerRepository) SaveCareer(ctx context.Context, userID, professionID uuid.UUID, notes string) error {
	query := `
		INSERT INTO saved_careers (user_id, profession_id, notes)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, profession_id)
		DO UPDATE SET notes = EXCLUDED.notes, saved_at = NOW()
	`

	_, err := r.db.Exec(ctx, query, userID, professionID, notes)
	if err != nil {
		return fmt.Errorf("failed to save career: %w", err)
	}

	return nil
}

// GetSavedCareers retrieves user's saved professions
func (r *CareerRepository) GetSavedCareers(ctx context.Context, userID uuid.UUID) ([]models.Profession, error) {
	query := `
		SELECT p.id, p.slug, p.title, p.category, p.ai_risk_score, p.avg_salary_usd,
		       p.description, p.is_future_job, p.demand_index, p.growth_pct, p.updated_at
		FROM professions p
		JOIN saved_careers sc ON p.id = sc.profession_id
		WHERE sc.user_id = $1
		ORDER BY sc.saved_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query saved careers: %w", err)
	}
	defer rows.Close()

	var professions []models.Profession
	for rows.Next() {
		var p models.Profession
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Category, &p.AIRiskScore,
			&p.AvgSalaryUSD, &p.Description, &p.IsFutureJob,
			&p.DemandIndex, &p.GrowthPct, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profession: %w", err)
		}
		professions = append(professions, p)
	}

	return professions, rows.Err()
}
