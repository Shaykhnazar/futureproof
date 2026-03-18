package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/internal/repository"
	"github.com/shaykhnazar/futureproof/pkg/cache"
)

const (
	citiesCacheKey = "cities:all"
	cityCacheTTL   = 1 * time.Hour
)

// CityService handles city-related business logic
type CityService struct {
	repo   *repository.CityRepository
	cache  *cache.Redis
	logger *zap.Logger
}

// NewCityService creates a new city service
func NewCityService(repo *repository.CityRepository, cache *cache.Redis, logger *zap.Logger) *CityService {
	return &CityService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

// GetAllCities retrieves all cities with caching
func (s *CityService) GetAllCities(ctx context.Context) ([]models.CityWithScore, error) {
	// Try cache first
	cached, err := s.cache.Get(ctx, citiesCacheKey)
	if err == nil && cached != "" {
		var cities []models.CityWithScore
		if err := json.Unmarshal([]byte(cached), &cities); err == nil {
			s.logger.Debug("Cities retrieved from cache")
			return cities, nil
		}
	}

	// Cache miss, query database
	cities, err := s.repo.GetAllCities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all cities: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(cities)
	_ = s.cache.Set(ctx, citiesCacheKey, data, cityCacheTTL)

	s.logger.Info("Retrieved cities from database", zap.Int("count", len(cities)))
	return cities, nil
}

// GetCityByID retrieves a city by ID
func (s *CityService) GetCityByID(ctx context.Context, id uuid.UUID) (*models.CityWithScore, error) {
	cacheKey := fmt.Sprintf("city:id:%s", id.String())

	// Try cache first
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var city models.CityWithScore
		if err := json.Unmarshal([]byte(cached), &city); err == nil {
			return &city, nil
		}
	}

	// Cache miss, query database
	city, err := s.repo.GetCityByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("city not found: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(city)
	_ = s.cache.Set(ctx, cacheKey, data, cityCacheTTL)

	return city, nil
}

// GetCitiesByRegion retrieves cities by region
func (s *CityService) GetCitiesByRegion(ctx context.Context, region string) ([]models.CityWithScore, error) {
	cacheKey := fmt.Sprintf("cities:region:%s", region)

	// Try cache first
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cities []models.CityWithScore
		if err := json.Unmarshal([]byte(cached), &cities); err == nil {
			return cities, nil
		}
	}

	// Cache miss, query database
	cities, err := s.repo.GetCitiesByRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities by region: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(cities)
	_ = s.cache.Set(ctx, cacheKey, data, cityCacheTTL)

	return cities, nil
}

// UpdateCityScore updates city opportunity score
func (s *CityService) UpdateCityScore(ctx context.Context, score models.CityScore) error {
	err := s.repo.UpdateCityScore(ctx, score)
	if err != nil {
		return fmt.Errorf("failed to update city score: %w", err)
	}

	// Invalidate cache
	_ = s.cache.Del(ctx, citiesCacheKey)
	_ = s.cache.Del(ctx, fmt.Sprintf("city:id:%s", score.CityID.String()))

	s.logger.Info("City score updated",
		zap.String("city_id", score.CityID.String()),
		zap.Int("score", score.Score),
	)

	return nil
}

// InvalidateCache invalidates city caches
func (s *CityService) InvalidateCache(ctx context.Context) error {
	_ = s.cache.Del(ctx, citiesCacheKey)
	s.logger.Info("City cache invalidated")
	return nil
}
