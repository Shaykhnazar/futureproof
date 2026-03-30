package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/internal/repository"
)

// numbeoCity maps city names to Numbeo query params and embedded fallback CLI.
// CLI = Cost of Living Index (NYC = 100 baseline).
// Source: Numbeo 2024 Cost of Living Index by City.
// Lower CLI = cheaper (better for talent).  We invert to a 0-100 score where
// 100 = very affordable and 0 = extremely expensive.
type numbeoCity struct {
	cityName    string
	country     string // Numbeo country name
	fallbackCLI float64
}

var numbeoCities = []numbeoCity{
	{"San Francisco", "United States", 93.5},
	{"New York", "United States", 100.0},
	{"Austin", "United States", 70.2},
	{"Seattle", "United States", 78.3},
	{"Toronto", "Canada", 65.4},
	{"London", "United Kingdom", 81.2},
	{"Berlin", "Germany", 64.8},
	{"Stockholm", "Sweden", 74.1},
	{"Amsterdam", "Netherlands", 79.3},
	{"Zurich", "Switzerland", 130.2},
	{"Tel Aviv", "Israel", 98.5},
	{"Dubai", "United Arab Emirates", 73.6},
	{"Singapore", "Singapore", 88.7},
	{"Tokyo", "Japan", 82.4},
	{"Seoul", "South Korea", 70.5},
	{"Bangalore", "India", 28.3},
	{"Sydney", "Australia", 84.1},
	{"Lagos", "Nigeria", 32.6},
	{"Nairobi", "Kenya", 35.2},
	{"Tashkent", "Uzbekistan", 24.8},
}

// NumbeoFetcher updates city cost_of_living scores from Numbeo.
// If no API key is configured it applies the embedded 2024 Numbeo CLI data.
type NumbeoFetcher struct {
	cityRepo   *repository.CityRepository
	logger     *zap.Logger
	apiKey     string
	httpClient *http.Client
}

func NewNumbeoFetcher(cityRepo *repository.CityRepository, logger *zap.Logger, apiKey string) *NumbeoFetcher {
	return &NumbeoFetcher{
		cityRepo:   cityRepo,
		logger:     logger,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

type numbeoAPIResponse struct {
	CostOfLivingIndex float64 `json:"cost_of_living_index"`
}

func (f *NumbeoFetcher) Run(ctx context.Context) error {
	f.logger.Info("Starting Numbeo cost-of-living update")

	cities, err := f.cityRepo.GetAllCities(ctx)
	if err != nil {
		return fmt.Errorf("failed to load cities: %w", err)
	}

	// Build lookup: city name → CityWithScore
	cityMap := make(map[string]models.CityWithScore, len(cities))
	for _, c := range cities {
		cityMap[c.Name] = c
	}

	for _, nc := range numbeoCities {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		city, ok := cityMap[nc.cityName]
		if !ok {
			f.logger.Debug("City not in DB, skipping", zap.String("city", nc.cityName))
			continue
		}

		cli := nc.fallbackCLI
		if f.apiKey != "" {
			if live, err := f.fetchCLI(ctx, nc.cityName, nc.country); err == nil {
				cli = live
			} else {
				f.logger.Warn("Numbeo live fetch failed, using fallback", zap.String("city", nc.cityName), zap.Error(err))
			}
		}

		// Convert CLI (NYC=100 baseline) to our inverted 0-100 score.
		// CLI 25 (very cheap) → score 88, CLI 130 (Zurich) → score 15
		colScore := int(100 - (cli / 130.0 * 85))
		if colScore < 0 {
			colScore = 0
		}
		if colScore > 100 {
			colScore = 100
		}

		score := models.CityScore{
			CityID:       city.ID,
			CostOfLiving: colScore,
			SnapshotDate: time.Now(),
			Source:       "numbeo",
		}

		if err := f.cityRepo.UpdateCityScore(ctx, score); err != nil {
			f.logger.Error("Failed to update city cost of living", zap.String("city", nc.cityName), zap.Error(err))
			continue
		}

		f.logger.Info("Updated cost of living",
			zap.String("city", nc.cityName),
			zap.Float64("cli", cli),
			zap.Int("score", colScore),
		)
	}

	f.logger.Info("Numbeo cost-of-living update complete")
	return nil
}

func (f *NumbeoFetcher) fetchCLI(ctx context.Context, city, country string) (float64, error) {
	params := url.Values{}
	params.Set("api_key", f.apiKey)
	params.Set("city", city)
	params.Set("country", country)
	params.Set("currency", "USD")

	reqURL := "https://www.numbeo.com/api/city_prices?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Numbeo returned %d", resp.StatusCode)
	}

	var result numbeoAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.CostOfLivingIndex == 0 {
		return 0, fmt.Errorf("zero CLI returned for %s", city)
	}

	return result.CostOfLivingIndex, nil
}
