package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/internal/repository"
)

// adzunaCountry maps country names to Adzuna country codes
var adzunaCountry = map[string]string{
	"USA":         "us",
	"UK":          "gb",
	"Germany":     "de",
	"Canada":      "ca",
	"Australia":   "au",
	"Singapore":   "sg",
}

type JobScraper struct {
	cityRepo   *repository.CityRepository
	logger     *zap.Logger
	apiKey     string
	appID      string
	httpClient *http.Client
}

func NewJobScraper(cityRepo *repository.CityRepository, logger *zap.Logger, apiKey, appID string) *JobScraper {
	return &JobScraper{
		cityRepo:   cityRepo,
		logger:     logger,
		apiKey:     apiKey,
		appID:      appID,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

type adzunaResponse struct {
	Count   int `json:"count"`
	Results []struct {
		Title    string  `json:"title"`
		SalaryMax float64 `json:"salary_max"`
	} `json:"results"`
}

// Run fetches tech job counts from Adzuna and uses them to update ai_investment scores
func (s *JobScraper) Run(ctx context.Context) error {
	if s.appID == "" || s.apiKey == "" || s.appID == "your-adzuna-app-id" {
		s.logger.Warn("Adzuna credentials not configured, skipping job scrape")
		return nil
	}

	s.logger.Info("Starting Adzuna job scrape")

	cities, err := s.cityRepo.GetAllCities(ctx)
	if err != nil {
		return fmt.Errorf("failed to load cities: %w", err)
	}

	for _, city := range cities {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		countryCode, ok := adzunaCountry[city.Country]
		if !ok {
			s.logger.Debug("No Adzuna country mapping, skipping", zap.String("country", city.Country))
			continue
		}

		// Fetch AI/tech job count for this city
		techJobs, err := s.fetchJobCount(ctx, countryCode, city.Name, "AI machine learning data science")
		if err != nil {
			s.logger.Warn("Adzuna fetch failed", zap.String("city", city.Name), zap.Error(err))
			continue
		}

		// Normalise: treat 1000+ tech jobs as score 100
		aiInvestment := int(float64(techJobs) / 10)
		if aiInvestment > 100 {
			aiInvestment = 100
		}

		score := models.CityScore{
			CityID:       city.ID,
			AIInvestment: aiInvestment,
			SnapshotDate: time.Now(),
			Source:       "adzuna",
		}

		if err := s.cityRepo.UpdateCityScore(ctx, score); err != nil {
			s.logger.Error("Failed to update city AI score", zap.String("city", city.Name), zap.Error(err))
			continue
		}

		s.logger.Info("Updated AI investment score",
			zap.String("city", city.Name),
			zap.Int("tech_jobs", techJobs),
			zap.Int("ai_investment", aiInvestment),
		)

		// Adzuna free tier: 1 req/s
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}

	s.logger.Info("Adzuna job scrape complete")
	return nil
}

func (s *JobScraper) fetchJobCount(ctx context.Context, country, city, keywords string) (int, error) {
	endpoint := fmt.Sprintf("https://api.adzuna.com/v1/api/jobs/%s/search/1", country)

	params := url.Values{}
	params.Set("app_id", s.appID)
	params.Set("app_key", s.apiKey)
	params.Set("where", city)
	params.Set("what", keywords)
	params.Set("results_per_page", "1")
	params.Set("content-type", "application/json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Adzuna returned status %d for %s/%s", resp.StatusCode, country, city)
	}

	var result adzunaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode Adzuna response: %w", err)
	}

	return result.Count, nil
}

// buildAIKeywords generates context-aware search terms
func buildAIKeywords(professions []string) string {
	keywords := []string{"AI", "machine learning", "data science", "software engineer"}
	for _, p := range professions {
		keywords = append(keywords, p)
	}
	return strings.Join(keywords[:4], " ")
}
