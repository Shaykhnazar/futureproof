package workers

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/repository"
)

// DataFetcher fetches economic data from external APIs (World Bank, BLS, etc.)
type DataFetcher struct {
	cityRepo       *repository.CityRepository
	logger         *zap.Logger
	worldBankURL   string
	blsAPIKey      string
}

// NewDataFetcher creates a new data fetcher
func NewDataFetcher(
	cityRepo *repository.CityRepository,
	logger *zap.Logger,
	worldBankURL string,
	blsAPIKey string,
) *DataFetcher {
	return &DataFetcher{
		cityRepo:     cityRepo,
		logger:       logger,
		worldBankURL: worldBankURL,
		blsAPIKey:    blsAPIKey,
	}
}

// Run executes the data fetching task
func (s *DataFetcher) Run(ctx context.Context) error {
	s.logger.Info("Starting economic data fetching")

	// Get all cities
	cities, err := s.cityRepo.GetAllCities(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch cities: %w", err)
	}

	s.logger.Info("Fetching economic data for cities", zap.Int("city_count", len(cities)))

	// TODO: Implement actual World Bank API / BLS API integration
	// For now, this is a placeholder

	for _, city := range cities {
		select {
		case <-ctx.Done():
			s.logger.Info("Data fetching cancelled")
			return ctx.Err()
		default:
			// Placeholder for actual API calls
			s.logger.Debug("Would fetch economic data for city",
				zap.String("city", city.Name),
				zap.String("country", city.Country),
			)

			// TODO: Fetch GDP growth, unemployment rate, etc.
			// data, err := s.fetchWorldBankData(city.Country)
			// if err != nil {
			//     s.logger.Error("Failed to fetch World Bank data", zap.Error(err))
			//     continue
			// }

			// TODO: Update city scores based on new data
			// s.updateCityScore(ctx, city.ID, data)
		}
	}

	s.logger.Info("Economic data fetching completed")
	return nil
}

// fetchWorldBankData would make the actual API call
// func (s *DataFetcher) fetchWorldBankData(country string) (map[string]interface{}, error) {
//     url := fmt.Sprintf("%s/country/%s/indicator/NY.GDP.MKTP.KD.ZG?format=json",
//         s.worldBankURL, country)
//
//     // Make HTTP request
//     // Parse response
//     // Return data
//     return nil, nil
// }
