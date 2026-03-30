package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/internal/repository"
)

// countryISO maps city countries to World Bank ISO-2 codes
var countryISO = map[string]string{
	"USA":         "US",
	"UK":          "GB",
	"Germany":     "DE",
	"Sweden":      "SE",
	"Netherlands": "NL",
	"Switzerland": "CH",
	"Singapore":   "SG",
	"Japan":       "JP",
	"South Korea": "KR",
	"India":       "IN",
	"Australia":   "AU",
	"Canada":      "CA",
	"Israel":      "IL",
	"UAE":         "AE",
	"Nigeria":     "NG",
	"Kenya":       "KE",
	"Uzbekistan":  "UZ",
}

type DataFetcher struct {
	cityRepo     *repository.CityRepository
	logger       *zap.Logger
	worldBankURL string
	httpClient   *http.Client
}

func NewDataFetcher(cityRepo *repository.CityRepository, logger *zap.Logger, worldBankURL, _ string) *DataFetcher {
	return &DataFetcher{
		cityRepo:     cityRepo,
		logger:       logger,
		worldBankURL: worldBankURL,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

type wbResponse struct {
	Value   *float64 `json:"value"`
	Country struct {
		Value string `json:"value"`
	} `json:"country"`
}

// Run fetches GDP growth & unemployment from the World Bank and updates city scores
func (s *DataFetcher) Run(ctx context.Context) error {
	s.logger.Info("Starting World Bank data fetch")

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

		iso, ok := countryISO[city.Country]
		if !ok {
			s.logger.Debug("No ISO mapping for country, skipping", zap.String("country", city.Country))
			continue
		}

		gdpGrowth, err := s.fetchIndicator(ctx, iso, "NY.GDP.MKTP.KD.ZG") // GDP growth %
		if err != nil {
			s.logger.Warn("Failed to fetch GDP growth", zap.String("city", city.Name), zap.Error(err))
			gdpGrowth = 0
		}

		unemployment, err := s.fetchIndicator(ctx, iso, "SL.UEM.TOTL.ZS") // Unemployment %
		if err != nil {
			s.logger.Warn("Failed to fetch unemployment", zap.String("city", city.Name), zap.Error(err))
			unemployment = 5.0
		}

		// Derive job_growth from GDP growth (rough proxy)
		jobGrowthPct := gdpGrowth * 0.6
		// Derive talent demand from inverse unemployment (lower unemployment → higher demand)
		talentDemand := int(100 - unemployment*5)
		if talentDemand < 0 {
			talentDemand = 0
		}
		if talentDemand > 100 {
			talentDemand = 100
		}

		score := models.CityScore{
			CityID:       city.ID,
			JobGrowthPct: jobGrowthPct,
			TalentDemand: talentDemand,
			SnapshotDate: time.Now(),
			Source:       "world_bank",
		}

		if err := s.cityRepo.UpdateCityScore(ctx, score); err != nil {
			s.logger.Error("Failed to update city score", zap.String("city", city.Name), zap.Error(err))
			continue
		}

		s.logger.Info("Updated city score from World Bank",
			zap.String("city", city.Name),
			zap.Float64("gdp_growth", gdpGrowth),
			zap.Float64("unemployment", unemployment),
		)
	}

	s.logger.Info("World Bank data fetch complete")
	return nil
}

func (s *DataFetcher) fetchIndicator(ctx context.Context, iso, indicator string) (float64, error) {
	url := fmt.Sprintf("%s/country/%s/indicator/%s?format=json&mrv=1&per_page=1", s.worldBankURL, iso, indicator)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// World Bank returns [{metadata}, [{data}]]
	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil || len(raw) < 2 {
		return 0, fmt.Errorf("unexpected World Bank response format")
	}

	var records []wbResponse
	if err := json.Unmarshal(raw[1], &records); err != nil || len(records) == 0 {
		return 0, fmt.Errorf("no records returned")
	}

	if records[0].Value == nil {
		return 0, fmt.Errorf("null value from World Bank")
	}
	return *records[0].Value, nil
}
