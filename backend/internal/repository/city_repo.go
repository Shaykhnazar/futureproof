package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/pkg/database"
)

// CityRepository handles city-related database operations
type CityRepository struct {
	db *database.DB
}

// NewCityRepository creates a new city repository
func NewCityRepository(db *database.DB) *CityRepository {
	return &CityRepository{db: db}
}

// GetAllCities retrieves all cities with their latest scores
func (r *CityRepository) GetAllCities(ctx context.Context) ([]models.CityWithScore, error) {
	query := `
		SELECT
			c.id, c.name, c.country, c.region, c.lat, c.lng,
			c.population, c.timezone, c.created_at, c.updated_at,
			COALESCE(cs.score, 0) as score,
			COALESCE(cs.job_growth_pct, 0) as job_growth_pct,
			COALESCE(cs.remote_score, 0) as remote_score,
			COALESCE(cs.ai_investment, 0) as ai_investment,
			COALESCE(cs.talent_demand, 0) as talent_demand,
			COALESCE(cs.cost_of_living, 0) as cost_of_living
		FROM cities c
		LEFT JOIN LATERAL (
			SELECT score, job_growth_pct, remote_score, ai_investment,
			       talent_demand, cost_of_living
			FROM city_scores
			WHERE city_id = c.id
			ORDER BY snapshot_date DESC
			LIMIT 1
		) cs ON true
		ORDER BY COALESCE(cs.score, 0) DESC, c.name ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query cities: %w", err)
	}
	defer rows.Close()

	var cities []models.CityWithScore
	for rows.Next() {
		var c models.CityWithScore
		err := rows.Scan(
			&c.ID, &c.Name, &c.Country, &c.Region, &c.Lat, &c.Lng,
			&c.Population, &c.Timezone, &c.CreatedAt, &c.UpdatedAt,
			&c.Score, &c.JobGrowthPct, &c.RemoteScore, &c.AIInvestment,
			&c.TalentDemand, &c.CostOfLiving,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan city: %w", err)
		}
		cities = append(cities, c)
	}

	return cities, rows.Err()
}

// GetCityByID retrieves a city by ID with its latest score
func (r *CityRepository) GetCityByID(ctx context.Context, id uuid.UUID) (*models.CityWithScore, error) {
	query := `
		SELECT
			c.id, c.name, c.country, c.region, c.lat, c.lng,
			c.population, c.timezone, c.created_at, c.updated_at,
			COALESCE(cs.score, 0) as score,
			COALESCE(cs.job_growth_pct, 0) as job_growth_pct,
			COALESCE(cs.remote_score, 0) as remote_score,
			COALESCE(cs.ai_investment, 0) as ai_investment,
			COALESCE(cs.talent_demand, 0) as talent_demand,
			COALESCE(cs.cost_of_living, 0) as cost_of_living
		FROM cities c
		LEFT JOIN LATERAL (
			SELECT score, job_growth_pct, remote_score, ai_investment,
			       talent_demand, cost_of_living
			FROM city_scores
			WHERE city_id = c.id
			ORDER BY snapshot_date DESC
			LIMIT 1
		) cs ON true
		WHERE c.id = $1
	`

	var c models.CityWithScore
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Country, &c.Region, &c.Lat, &c.Lng,
		&c.Population, &c.Timezone, &c.CreatedAt, &c.UpdatedAt,
		&c.Score, &c.JobGrowthPct, &c.RemoteScore, &c.AIInvestment,
		&c.TalentDemand, &c.CostOfLiving,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get city: %w", err)
	}

	return &c, nil
}

// GetCitiesByRegion retrieves cities filtered by region
func (r *CityRepository) GetCitiesByRegion(ctx context.Context, region string) ([]models.CityWithScore, error) {
	query := `
		SELECT
			c.id, c.name, c.country, c.region, c.lat, c.lng,
			c.population, c.timezone, c.created_at, c.updated_at,
			COALESCE(cs.score, 0) as score,
			COALESCE(cs.job_growth_pct, 0) as job_growth_pct,
			COALESCE(cs.remote_score, 0) as remote_score,
			COALESCE(cs.ai_investment, 0) as ai_investment,
			COALESCE(cs.talent_demand, 0) as talent_demand,
			COALESCE(cs.cost_of_living, 0) as cost_of_living
		FROM cities c
		LEFT JOIN LATERAL (
			SELECT score, job_growth_pct, remote_score, ai_investment,
			       talent_demand, cost_of_living
			FROM city_scores
			WHERE city_id = c.id
			ORDER BY snapshot_date DESC
			LIMIT 1
		) cs ON true
		WHERE c.region = $1
		ORDER BY COALESCE(cs.score, 0) DESC
	`

	rows, err := r.db.Query(ctx, query, region)
	if err != nil {
		return nil, fmt.Errorf("failed to query cities by region: %w", err)
	}
	defer rows.Close()

	var cities []models.CityWithScore
	for rows.Next() {
		var c models.CityWithScore
		err := rows.Scan(
			&c.ID, &c.Name, &c.Country, &c.Region, &c.Lat, &c.Lng,
			&c.Population, &c.Timezone, &c.CreatedAt, &c.UpdatedAt,
			&c.Score, &c.JobGrowthPct, &c.RemoteScore, &c.AIInvestment,
			&c.TalentDemand, &c.CostOfLiving,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan city: %w", err)
		}
		cities = append(cities, c)
	}

	return cities, rows.Err()
}

// UpdateCityScore updates or creates a city score
func (r *CityRepository) UpdateCityScore(ctx context.Context, score models.CityScore) error {
	query := `
		INSERT INTO city_scores (
			city_id, score, job_growth_pct, remote_score,
			ai_investment, talent_demand, cost_of_living, source
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (city_id, snapshot_date)
		DO UPDATE SET
			score = EXCLUDED.score,
			job_growth_pct = EXCLUDED.job_growth_pct,
			remote_score = EXCLUDED.remote_score,
			ai_investment = EXCLUDED.ai_investment,
			talent_demand = EXCLUDED.talent_demand,
			cost_of_living = EXCLUDED.cost_of_living,
			source = EXCLUDED.source
	`

	_, err := r.db.Exec(ctx, query,
		score.CityID, score.Score, score.JobGrowthPct, score.RemoteScore,
		score.AIInvestment, score.TalentDemand, score.CostOfLiving, score.Source,
	)
	if err != nil {
		return fmt.Errorf("failed to update city score: %w", err)
	}

	return nil
}
