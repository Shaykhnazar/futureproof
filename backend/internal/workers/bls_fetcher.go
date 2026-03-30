package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/repository"
)

// blsProfession maps a slug to its BLS SOC code and embedded official data.
// Growth is BLS Employment Projections 2022–2032 (published).
// AIRisk is the Frey & Osborne (2013/2017) automation probability (0–100).
type blsProfession struct {
	slug       string
	socCode    string  // 6-digit SOC without hyphen, e.g. "151252"
	growth2232 float64 // BLS 10-year projected employment change %
	aiRisk     int     // Frey & Osborne automation probability
}

// blsOccupations is the authoritative reference table.
// Sources:
//   - Salaries: BLS OES national estimates (fetched live via API)
//   - Growth:   BLS Employment Projections 2022–2032
//   - AI Risk:  Frey & Osborne (2013), updated McKinsey 2023 mapping
var blsOccupations = []blsProfession{
	{"software-engineer", "151252", 26.0, 45},
	{"data-scientist", "152051", 35.0, 38},
	{"cybersecurity-analyst", "151212", 32.0, 28},
	{"devops-engineer", "151244", 5.0, 42},
	{"ux-designer", "151255", 3.0, 40},
	{"graphic-designer", "271024", -3.0, 75},
	{"content-writer", "273043", -4.0, 80},
	{"financial-analyst", "132051", 8.0, 85},
	{"accountant", "132011", -4.0, 94},
	{"marketing-manager", "112021", 6.0, 60},
	{"project-manager", "119198", 7.0, 50},
	{"lawyer", "231011", 8.0, 55},
	{"doctor", "291216", 3.0, 25},
	{"nurse", "291141", 6.0, 20},
	{"teacher", "252031", 1.0, 30},
	{"truck-driver", "533032", 4.0, 90},
}

// BLSFetcher fetches real salary data from BLS OES and applies embedded
// Frey & Osborne AI risk scores and BLS 2022-2032 employment projections.
type BLSFetcher struct {
	careerRepo *repository.CareerRepository
	logger     *zap.Logger
	apiKey     string
	httpClient *http.Client
}

func NewBLSFetcher(careerRepo *repository.CareerRepository, logger *zap.Logger, apiKey string) *BLSFetcher {
	return &BLSFetcher{
		careerRepo: careerRepo,
		logger:     logger,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (f *BLSFetcher) Run(ctx context.Context) error {
	f.logger.Info("Starting BLS profession data update")

	// Build OES series IDs for mean annual wage (data type 03, national area 0000000)
	seriesIDs := make([]string, 0, len(blsOccupations))
	socToSlug := make(map[string]string)
	for _, p := range blsOccupations {
		sid := "OES0000000" + p.socCode + "03"
		seriesIDs = append(seriesIDs, sid)
		socToSlug[sid] = p.slug
	}

	wages, err := f.fetchWages(ctx, seriesIDs)
	if err != nil {
		f.logger.Warn("BLS API fetch failed, applying embedded projections with $0 salary update", zap.Error(err))
	}

	for _, p := range blsOccupations {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		sid := "OES0000000" + p.socCode + "03"
		salary := wages[sid] // 0 if fetch failed — UpdateProfessionStats skips 0-salary update

		// demand_index: 50 baseline shifted by growth (capped 0–100)
		demandIndex := 50 + int(p.growth2232*2)
		if demandIndex < 0 {
			demandIndex = 0
		}
		if demandIndex > 100 {
			demandIndex = 100
		}

		if err := f.careerRepo.UpdateProfessionStats(ctx, p.slug, salary, p.growth2232, demandIndex, p.aiRisk); err != nil {
			f.logger.Error("Failed to update profession stats", zap.String("slug", p.slug), zap.Error(err))
			continue
		}

		f.logger.Info("Updated profession",
			zap.String("slug", p.slug),
			zap.Int("salary_usd", salary),
			zap.Float64("growth_pct", p.growth2232),
			zap.Int("ai_risk", p.aiRisk),
		)
	}

	f.logger.Info("BLS profession update complete")
	return nil
}

type blsRequest struct {
	SeriesID        []string `json:"seriesid"`
	StartYear       string   `json:"startyear"`
	EndYear         string   `json:"endyear"`
	RegistrationKey string   `json:"registrationkey,omitempty"`
}

type blsAPIResponse struct {
	Status  string `json:"status"`
	Message []string `json:"message"`
	Results struct {
		Series []struct {
			SeriesID string `json:"seriesId"`
			Data     []struct {
				Year  string `json:"year"`
				Value string `json:"value"`
			} `json:"data"`
		} `json:"series"`
	} `json:"Results"`
}

func (f *BLSFetcher) fetchWages(ctx context.Context, seriesIDs []string) (map[string]int, error) {
	year := strconv.Itoa(time.Now().Year() - 1) // OES lags ~1 year

	reqBody := blsRequest{
		SeriesID:  seriesIDs,
		StartYear: year,
		EndYear:   year,
	}
	if f.apiKey != "" {
		reqBody.RegistrationKey = f.apiKey
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.bls.gov/publicAPI/v2/timeseries/data/",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("BLS request failed: %w", err)
	}
	defer resp.Body.Close()

	var result blsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("BLS decode failed: %w", err)
	}

	if result.Status != "REQUEST_SUCCEEDED" {
		return nil, fmt.Errorf("BLS API error: %s", strings.Join(result.Message, "; "))
	}

	wages := make(map[string]int, len(result.Results.Series))
	for _, s := range result.Results.Series {
		if len(s.Data) == 0 {
			continue
		}
		// BLS returns wage as string with commas, e.g. "145,230"
		cleaned := strings.ReplaceAll(s.Data[0].Value, ",", "")
		if w, err := strconv.Atoi(cleaned); err == nil {
			wages[s.SeriesID] = w
		}
	}

	return wages, nil
}
