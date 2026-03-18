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
	professionsCacheKey = "professions:all"
	professionCacheTTL  = 1 * time.Hour
	transitionsCacheTTL = 1 * time.Hour
)

// CareerService handles career-related business logic
type CareerService struct {
	repo   *repository.CareerRepository
	cache  *cache.Redis
	logger *zap.Logger
}

// NewCareerService creates a new career service
func NewCareerService(repo *repository.CareerRepository, cache *cache.Redis, logger *zap.Logger) *CareerService {
	return &CareerService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

// GetAllProfessions retrieves all professions with caching
func (s *CareerService) GetAllProfessions(ctx context.Context) ([]models.Profession, error) {
	// Try cache first
	cached, err := s.cache.Get(ctx, professionsCacheKey)
	if err == nil && cached != "" {
		var professions []models.Profession
		if err := json.Unmarshal([]byte(cached), &professions); err == nil {
			s.logger.Debug("Professions retrieved from cache")
			return professions, nil
		}
	}

	// Cache miss, query database
	professions, err := s.repo.GetAllProfessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all professions: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(professions)
	_ = s.cache.Set(ctx, professionsCacheKey, data, professionCacheTTL)

	s.logger.Info("Retrieved professions from database", zap.Int("count", len(professions)))
	return professions, nil
}

// GetProfessionBySlug retrieves a profession by slug with caching
func (s *CareerService) GetProfessionBySlug(ctx context.Context, slug string) (*models.Profession, error) {
	cacheKey := fmt.Sprintf("profession:slug:%s", slug)

	// Try cache first
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var profession models.Profession
		if err := json.Unmarshal([]byte(cached), &profession); err == nil {
			return &profession, nil
		}
	}

	// Cache miss, query database
	profession, err := s.repo.GetProfessionBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("profession not found: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(profession)
	_ = s.cache.Set(ctx, cacheKey, data, professionCacheTTL)

	return profession, nil
}

// GetCareerTransitions retrieves career pivot recommendations
func (s *CareerService) GetCareerTransitions(ctx context.Context, professionSlug string) ([]models.CareerTransitionWithDetails, error) {
	cacheKey := fmt.Sprintf("transitions:%s", professionSlug)

	// Try cache first
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var transitions []models.CareerTransitionWithDetails
		if err := json.Unmarshal([]byte(cached), &transitions); err == nil {
			return transitions, nil
		}
	}

	// Cache miss, query database
	transitions, err := s.repo.GetCareerTransitions(ctx, professionSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to get career transitions: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(transitions)
	_ = s.cache.Set(ctx, cacheKey, data, transitionsCacheTTL)

	s.logger.Info("Retrieved career transitions",
		zap.String("profession", professionSlug),
		zap.Int("count", len(transitions)),
	)
	return transitions, nil
}

// GetFutureProfessions retrieves all future/emerging professions
func (s *CareerService) GetFutureProfessions(ctx context.Context) ([]models.Profession, error) {
	cacheKey := "professions:future"

	// Try cache first
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var professions []models.Profession
		if err := json.Unmarshal([]byte(cached), &professions); err == nil {
			return professions, nil
		}
	}

	// Cache miss, query database
	professions, err := s.repo.GetFutureProfessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get future professions: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(professions)
	_ = s.cache.Set(ctx, cacheKey, data, professionCacheTTL)

	return professions, nil
}

// SaveCareer saves a profession to user's list
func (s *CareerService) SaveCareer(ctx context.Context, userID, professionID uuid.UUID, notes string) error {
	err := s.repo.SaveCareer(ctx, userID, professionID, notes)
	if err != nil {
		return fmt.Errorf("failed to save career: %w", err)
	}

	s.logger.Info("Career saved",
		zap.String("user_id", userID.String()),
		zap.String("profession_id", professionID.String()),
	)
	return nil
}

// GetSavedCareers retrieves user's saved professions
func (s *CareerService) GetSavedCareers(ctx context.Context, userID uuid.UUID) ([]models.Profession, error) {
	professions, err := s.repo.GetSavedCareers(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get saved careers: %w", err)
	}

	return professions, nil
}

// InvalidateCache invalidates profession caches
func (s *CareerService) InvalidateCache(ctx context.Context) error {
	keys := []string{
		professionsCacheKey,
		"professions:future",
	}

	for _, key := range keys {
		_ = s.cache.Del(ctx, key)
	}

	s.logger.Info("Profession cache invalidated")
	return nil
}
