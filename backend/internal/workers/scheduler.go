package workers

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/config"
)

// Scheduler manages periodic background tasks
type Scheduler struct {
	config         *config.WorkersConfig
	logger         *zap.Logger
	jobScraper     *JobScraper
	dataFetcher    *DataFetcher
	blsFetcher     *BLSFetcher
	numbeoFetcher  *NumbeoFetcher
	stopChan       chan struct{}
}

// NewScheduler creates a new task scheduler
func NewScheduler(
	cfg *config.WorkersConfig,
	logger *zap.Logger,
	jobScraper *JobScraper,
	dataFetcher *DataFetcher,
	blsFetcher *BLSFetcher,
	numbeoFetcher *NumbeoFetcher,
) *Scheduler {
	return &Scheduler{
		config:        cfg,
		logger:        logger,
		jobScraper:    jobScraper,
		dataFetcher:   dataFetcher,
		blsFetcher:    blsFetcher,
		numbeoFetcher: numbeoFetcher,
		stopChan:      make(chan struct{}),
	}
}

// Start begins running scheduled tasks
func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info("Starting background task scheduler")

	go s.runPeriodic(ctx, "Job Scraper", s.config.ScraperInterval, s.jobScraper.Run)
	go s.runPeriodic(ctx, "World Bank Fetcher", s.config.DataFetchInterval, s.dataFetcher.Run)
	go s.runPeriodic(ctx, "BLS Fetcher", s.config.BLSFetchInterval, s.blsFetcher.Run)
	go s.runPeriodic(ctx, "Numbeo Fetcher", s.config.BLSFetchInterval, s.numbeoFetcher.Run)
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	s.logger.Info("Stopping background task scheduler")
	close(s.stopChan)
}

// runPeriodic runs a task periodically
func (s *Scheduler) runPeriodic(ctx context.Context, name string, interval time.Duration, task func(context.Context) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately on startup
	s.logger.Info("Running initial task", zap.String("task", name))
	if err := task(ctx); err != nil {
		s.logger.Error("Initial task failed", zap.String("task", name), zap.Error(err))
	}

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Running scheduled task", zap.String("task", name))
			if err := task(ctx); err != nil {
				s.logger.Error("Scheduled task failed", zap.String("task", name), zap.Error(err))
			}

		case <-s.stopChan:
			s.logger.Info("Task stopped", zap.String("task", name))
			return

		case <-ctx.Done():
			s.logger.Info("Context cancelled, stopping task", zap.String("task", name))
			return
		}
	}
}
