package workers

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/config"
)

// Scheduler manages periodic background tasks
type Scheduler struct {
	config      *config.WorkersConfig
	logger      *zap.Logger
	jobScraper  *JobScraper
	dataFetcher *DataFetcher
	stopChan    chan struct{}
}

// NewScheduler creates a new task scheduler
func NewScheduler(
	cfg *config.WorkersConfig,
	logger *zap.Logger,
	jobScraper *JobScraper,
	dataFetcher *DataFetcher,
) *Scheduler {
	return &Scheduler{
		config:      cfg,
		logger:      logger,
		jobScraper:  jobScraper,
		dataFetcher: dataFetcher,
		stopChan:    make(chan struct{}),
	}
}

// Start begins running scheduled tasks
func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info("Starting background task scheduler")

	// Start job scraper
	go s.runPeriodic(ctx, "Job Scraper", s.config.ScraperInterval, s.jobScraper.Run)

	// Start data fetcher
	go s.runPeriodic(ctx, "Data Fetcher", s.config.DataFetchInterval, s.dataFetcher.Run)
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
