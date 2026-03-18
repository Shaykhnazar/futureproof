// ============================================================
// FILE: go.mod
// ============================================================
//
// module github.com/yourname/futureproof-api
//
// go 1.22
//
// require (
//     github.com/gofiber/fiber/v2       v2.52.4
//     github.com/gofiber/websocket/v2   v2.2.1
//     github.com/golang-jwt/jwt/v5      v5.2.1
//     github.com/jackc/pgx/v5           v5.6.0
//     github.com/redis/go-redis/v9      v9.5.1
//     github.com/joho/godotenv          v1.5.1
//     github.com/anthropics/anthropic-sdk-go v0.2.0
//     github.com/robfig/cron/v3         v3.0.1
//     go.uber.org/zap                   v1.27.0
//     golang.org/x/crypto               v0.22.0
// )

// ============================================================
// FILE: cmd/server/main.go
// ============================================================

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yourname/futureproof-api/internal/api"
	"github.com/yourname/futureproof-api/internal/config"
	"github.com/yourname/futureproof-api/internal/workers"
	"github.com/yourname/futureproof-api/pkg/cache"
	"github.com/yourname/futureproof-api/pkg/database"
	"github.com/yourname/futureproof-api/pkg/logger"
)

func main() {
	// Load config from .env
	cfg := config.Load()
	log := logger.New(cfg.Env)

	// Connect to PostgreSQL
	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("postgres connection failed", zap.Error(err))
	}
	defer db.Close()

	// Connect to Redis
	rdb, err := cache.NewRedis(cfg.RedisURL)
	if err != nil {
		log.Fatal("redis connection failed", zap.Error(err))
	}
	defer rdb.Close()

	// Build Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "FutureProof API v1",
		ErrorHandler: api.ErrorHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Register all routes
	api.SetupRoutes(app, db, rdb, cfg, log)

	// Start background workers
	scheduler := workers.NewScheduler(db, rdb, cfg, log)
	scheduler.Start()
	defer scheduler.Stop()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server starting", zap.String("port", cfg.Port))
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Error("server error", zap.Error(err))
		}
	}()

	<-quit
	log.Info("shutting down gracefully…")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(ctx)
}

// ============================================================
// FILE: internal/config/config.go
// ============================================================

package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env            string
	Port           string
	DatabaseURL    string
	RedisURL       string
	AnthropicKey   string
	JWTSecret      string
	AllowedOrigins string
	AdzunaAppID    string
	AdzunaKey      string
}

func Load() *Config {
	_ = godotenv.Load() // no-op in production (env vars already set)
	return &Config{
		Env:            getEnv("APP_ENV", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://fp:fp@localhost:5432/futureproof"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		AnthropicKey:   mustEnv("ANTHROPIC_API_KEY"),
		JWTSecret:      mustEnv("JWT_SECRET"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		AdzunaAppID:    getEnv("ADZUNA_APP_ID", ""),
		AdzunaKey:      getEnv("ADZUNA_KEY", ""),
	}
}

func getEnv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("required env var missing: " + k)
	}
	return v
}

// ============================================================
// FILE: pkg/database/postgres.go
// ============================================================

package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewPostgres(url string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns           = 20
	cfg.MinConns           = 2
	cfg.MaxConnLifetime    = 30 * time.Minute
	cfg.MaxConnIdleTime    = 5 * time.Minute
	cfg.HealthCheckPeriod  = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() { db.Pool.Close() }

// ============================================================
// FILE: pkg/cache/redis.go
// ============================================================

package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewRedis(url string) (*Cache, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	c := redis.NewClient(opt)
	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Cache{client: c}, nil
}

func (c *Cache) Close() error { return c.client.Close() }

func (c *Cache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, b, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	b, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err // redis.Nil if not found
	}
	return json.Unmarshal(b, dest)
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// ============================================================
// FILE: internal/models/city.go
// ============================================================

package models

import "time"

type City struct {
	ID         string    `db:"id"          json:"id"`
	Name       string    `db:"name"        json:"name"`
	Country    string    `db:"country"     json:"country"`
	Region     string    `db:"region"      json:"region"`
	Lat        float64   `db:"lat"         json:"lat"`
	Lng        float64   `db:"lng"         json:"lng"`
	Population int       `db:"population"  json:"population,omitempty"`
	Timezone   string    `db:"timezone"    json:"timezone,omitempty"`
	CreatedAt  time.Time `db:"created_at"  json:"-"`
}

type CityScore struct {
	ID           string    `db:"id"             json:"id"`
	CityID       string    `db:"city_id"        json:"city_id"`
	Score        int       `db:"score"          json:"score"`
	JobGrowthPct float64   `db:"job_growth_pct" json:"job_growth_pct"`
	RemoteScore  int       `db:"remote_score"   json:"remote_score"`
	AIInvestment int       `db:"ai_investment"  json:"ai_investment"`
	TalentDemand int       `db:"talent_demand"  json:"talent_demand"`
	CostOfLiving int       `db:"cost_of_living" json:"cost_of_living"`
	SnapshotDate time.Time `db:"snapshot_date"  json:"snapshot_date"`
	Source       string    `db:"source"         json:"source"`
}

// CityDetail combines city + latest score + top professions
type CityDetail struct {
	City
	Score        int      `json:"score"`
	JobGrowthPct float64  `json:"job_growth_pct"`
	RemoteScore  int      `json:"remote_score"`
	AIInvestment int      `json:"ai_investment"`
	TopJobs      []string `json:"top_jobs"`
}

// ============================================================
// FILE: internal/models/profession.go
// ============================================================

package models

import "time"

type Profession struct {
	ID           string    `db:"id"             json:"id"`
	Slug         string    `db:"slug"           json:"slug"`
	Title        string    `db:"title"          json:"title"`
	Category     string    `db:"category"       json:"category"`
	AIRiskScore  int       `db:"ai_risk_score"  json:"ai_risk_score"`
	AvgSalaryUSD int       `db:"avg_salary_usd" json:"avg_salary_usd"`
	Description  string    `db:"description"    json:"description"`
	IsFutureJob  bool      `db:"is_future_job"  json:"is_future_job"`
	DemandIndex  int       `db:"demand_index"   json:"demand_index"`
	GrowthPct    float64   `db:"growth_pct"     json:"growth_pct"`
	UpdatedAt    time.Time `db:"updated_at"     json:"-"`
}

type Skill struct {
	ID        string `db:"id"          json:"id"`
	Name      string `db:"name"         json:"name"`
	Category  string `db:"category"     json:"category"`
	IsAIProof bool   `db:"is_ai_proof"  json:"is_ai_proof"`
}

type CareerTransition struct {
	FromProfession    string `db:"from_profession"    json:"from_profession"`
	ToProfession      string `db:"to_profession"      json:"to_profession"`
	ToProfessionTitle string `db:"to_profession_title" json:"to_profession_title"`
	MatchScore        int    `db:"match_score"        json:"match_score"`
	TransitionReason  string `db:"transition_reason"  json:"transition_reason"`
	AvgReskillMonths  int    `db:"avg_reskill_months" json:"avg_reskill_months"`
}

// ============================================================
// FILE: internal/models/user.go
// ============================================================

package models

import "time"

type User struct {
	ID           string    `db:"id"            json:"id"`
	Email        string    `db:"email"         json:"email"`
	Name         string    `db:"name"          json:"name"`
	AvatarURL    string    `db:"avatar_url"    json:"avatar_url,omitempty"`
	AuthProvider string    `db:"auth_provider" json:"auth_provider"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}

type UserProfile struct {
	UserID       string   `db:"user_id"        json:"user_id"`
	CurrentJobID string   `db:"current_job_id" json:"current_job_id,omitempty"`
	CityID       string   `db:"city_id"        json:"city_id,omitempty"`
	YearsExp     int      `db:"years_exp"      json:"years_exp"`
	Education    string   `db:"education"      json:"education"`
	TargetJobID  string   `db:"target_job_id"  json:"target_job_id,omitempty"`
	Skills       []string `db:"skills"         json:"skills"`
}

// ============================================================
// FILE: internal/models/analysis.go
// ============================================================

package models

import "time"

type AnalysisRequest struct {
	Profession string   `json:"profession" validate:"required"`
	Location   string   `json:"location"`
	YearsExp   int      `json:"years_experience"`
	Skills     []string `json:"skills"`
}

type CareerAnalysis struct {
	Profession       string             `json:"profession"`
	AIRiskScore      int                `json:"ai_risk_score"`
	RiskLevel        string             `json:"risk_level"`
	Explanation      string             `json:"explanation"`
	ResilientSkills  []string           `json:"resilient_skills"`
	AtRiskTasks      []string           `json:"at_risk_tasks"`
	TopPivots        []PivotSuggestion  `json:"top_pivots"`
	Timeline         string             `json:"timeline"`
	ImmediateActions []string           `json:"immediate_actions"`
	CachedAt         time.Time          `json:"cached_at"`
}

type PivotSuggestion struct {
	Title           string   `json:"title"`
	MatchScore      int      `json:"match_score"`
	SkillGap        []string `json:"skill_gap"`
	AvgReskillMonths int     `json:"avg_reskill_months"`
	Reasoning       string   `json:"reasoning"`
	AvgSalaryUSD    int      `json:"avg_salary_usd"`
	GrowthPct       string   `json:"growth_pct"`
}

// ============================================================
// FILE: internal/repository/city_repo.go
// ============================================================

package repository

import (
	"context"
	"fmt"

	"github.com/yourname/futureproof-api/internal/models"
	"github.com/yourname/futureproof-api/pkg/database"
)

type CityRepo struct{ db *database.DB }

func NewCityRepo(db *database.DB) *CityRepo { return &CityRepo{db} }

func (r *CityRepo) GetAll(ctx context.Context) ([]models.CityDetail, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT
			c.id, c.name, c.country, c.region, c.lat, c.lng,
			cs.score, cs.job_growth_pct, cs.remote_score, cs.ai_investment
		FROM cities c
		JOIN LATERAL (
			SELECT * FROM city_scores
			WHERE city_id = c.id
			ORDER BY snapshot_date DESC LIMIT 1
		) cs ON true
		ORDER BY cs.score DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("city_repo.GetAll: %w", err)
	}
	defer rows.Close()

	var cities []models.CityDetail
	for rows.Next() {
		var c models.CityDetail
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Country, &c.Region, &c.Lat, &c.Lng,
			&c.Score, &c.JobGrowthPct, &c.RemoteScore, &c.AIInvestment,
		); err != nil {
			return nil, err
		}
		cities = append(cities, c)
	}
	return cities, nil
}

func (r *CityRepo) GetByID(ctx context.Context, id string) (*models.CityDetail, error) {
	var c models.CityDetail
	err := r.db.Pool.QueryRow(ctx, `
		SELECT
			c.id, c.name, c.country, c.region, c.lat, c.lng, c.population, c.timezone,
			cs.score, cs.job_growth_pct, cs.remote_score, cs.ai_investment, cs.talent_demand
		FROM cities c
		JOIN LATERAL (
			SELECT * FROM city_scores WHERE city_id = c.id ORDER BY snapshot_date DESC LIMIT 1
		) cs ON true
		WHERE c.id = $1
	`, id).Scan(
		&c.ID, &c.Name, &c.Country, &c.Region, &c.Lat, &c.Lng, &c.Population, &c.Timezone,
		&c.Score, &c.JobGrowthPct, &c.RemoteScore, &c.AIInvestment,
	)
	if err != nil {
		return nil, fmt.Errorf("city_repo.GetByID(%s): %w", id, err)
	}
	// Fetch top jobs
	jobRows, _ := r.db.Pool.Query(ctx, `
		SELECT p.title FROM city_top_professions ctp
		JOIN professions p ON p.id = ctp.profession_id
		WHERE ctp.city_id = $1
		ORDER BY ctp.rank ASC LIMIT 5
	`, id)
	defer jobRows.Close()
	for jobRows.Next() {
		var t string
		_ = jobRows.Scan(&t)
		c.TopJobs = append(c.TopJobs, t)
	}
	return &c, nil
}

func (r *CityRepo) Search(ctx context.Context, q string) ([]models.CityDetail, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT c.id, c.name, c.country, c.region, c.lat, c.lng, cs.score
		FROM cities c
		JOIN LATERAL (
			SELECT score FROM city_scores WHERE city_id = c.id ORDER BY snapshot_date DESC LIMIT 1
		) cs ON true
		WHERE c.name ILIKE $1 OR c.country ILIKE $1
		ORDER BY cs.score DESC LIMIT 15
	`, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.CityDetail
	for rows.Next() {
		var c models.CityDetail
		_ = rows.Scan(&c.ID, &c.Name, &c.Country, &c.Region, &c.Lat, &c.Lng, &c.Score)
		out = append(out, c)
	}
	return out, nil
}

// ============================================================
// FILE: internal/repository/career_repo.go
// ============================================================

package repository

import (
	"context"
	"fmt"

	"github.com/yourname/futureproof-api/internal/models"
	"github.com/yourname/futureproof-api/pkg/database"
)

type CareerRepo struct{ db *database.DB }

func NewCareerRepo(db *database.DB) *CareerRepo { return &CareerRepo{db} }

func (r *CareerRepo) GetAll(ctx context.Context) ([]models.Profession, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, slug, title, category, ai_risk_score, avg_salary_usd,
		       description, is_future_job, demand_index, growth_pct
		FROM professions ORDER BY title ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("career_repo.GetAll: %w", err)
	}
	defer rows.Close()
	var out []models.Profession
	for rows.Next() {
		var p models.Profession
		_ = rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Category, &p.AIRiskScore,
			&p.AvgSalaryUSD, &p.Description, &p.IsFutureJob, &p.DemandIndex, &p.GrowthPct)
		out = append(out, p)
	}
	return out, nil
}

func (r *CareerRepo) GetBySlug(ctx context.Context, slug string) (*models.Profession, error) {
	var p models.Profession
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, slug, title, category, ai_risk_score, avg_salary_usd,
		       description, is_future_job, demand_index, growth_pct
		FROM professions WHERE slug = $1
	`, slug).Scan(&p.ID, &p.Slug, &p.Title, &p.Category, &p.AIRiskScore,
		&p.AvgSalaryUSD, &p.Description, &p.IsFutureJob, &p.DemandIndex, &p.GrowthPct)
	if err != nil {
		return nil, fmt.Errorf("career_repo.GetBySlug(%s): %w", slug, err)
	}
	return &p, nil
}

func (r *CareerRepo) GetTransitions(ctx context.Context, fromSlug string) ([]models.CareerTransition, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT ct.from_profession, ct.to_profession, p.title,
		       ct.match_score, ct.transition_reason, ct.avg_reskill_months
		FROM career_transitions ct
		JOIN professions p ON p.id = ct.to_profession
		JOIN professions fp ON fp.id = ct.from_profession AND fp.slug = $1
		ORDER BY ct.match_score DESC LIMIT 6
	`, fromSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.CareerTransition
	for rows.Next() {
		var t models.CareerTransition
		_ = rows.Scan(&t.FromProfession, &t.ToProfession, &t.ToProfessionTitle,
			&t.MatchScore, &t.TransitionReason, &t.AvgReskillMonths)
		out = append(out, t)
	}
	return out, nil
}

func (r *CareerRepo) GetFutureJobs(ctx context.Context) ([]models.Profession, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, slug, title, category, ai_risk_score, avg_salary_usd,
		       description, is_future_job, demand_index, growth_pct
		FROM professions WHERE is_future_job = true
		ORDER BY demand_index DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Profession
	for rows.Next() {
		var p models.Profession
		_ = rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Category, &p.AIRiskScore,
			&p.AvgSalaryUSD, &p.Description, &p.IsFutureJob, &p.DemandIndex, &p.GrowthPct)
		out = append(out, p)
	}
	return out, nil
}

// ============================================================
// FILE: internal/services/ai_service.go
// ============================================================

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/yourname/futureproof-api/internal/config"
	"github.com/yourname/futureproof-api/internal/models"
	"github.com/yourname/futureproof-api/pkg/cache"
)

type AIService struct {
	client *anthropic.Client
	cache  *cache.Cache
	cfg    *config.Config
}

func NewAIService(cfg *config.Config, c *cache.Cache) *AIService {
	client := anthropic.NewClient(
		anthropic.WithAPIKey(cfg.AnthropicKey),
	)
	return &AIService{client: client, cache: c, cfg: cfg}
}

// AnalyzeCareer returns a full career risk analysis from Claude (cached 24h)
func (s *AIService) AnalyzeCareer(ctx context.Context, req models.AnalysisRequest) (*models.CareerAnalysis, error) {
	cacheKey := fmt.Sprintf("analysis:%s:%s:%d", slugify(req.Profession), slugify(req.Location), req.YearsExp)

	// Cache check
	var cached models.CareerAnalysis
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	prompt := buildAnalysisPrompt(req)

	msg, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(2048)),
		System: anthropic.F([]anthropic.TextBlockParam{
			{Text: anthropic.F(systemPromptAnalysis)},
		}),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.UserMessageParam(anthropic.NewTextBlock(prompt)),
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("anthropic API error: %w", err)
	}

	raw := msg.Content[0].Text
	// Claude may wrap JSON in ```json ... ``` — strip fences
	raw = strings.TrimPrefix(strings.TrimSpace(raw), "```json")
	raw = strings.TrimSuffix(raw, "```")

	var analysis models.CareerAnalysis
	if err := json.Unmarshal([]byte(raw), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}
	analysis.CachedAt = time.Now()

	// Cache 24 hours
	_ = s.cache.Set(ctx, cacheKey, analysis, 24*time.Hour)
	return &analysis, nil
}

// StreamChat sends a multi-turn career coach message and streams the response via a channel
func (s *AIService) StreamChat(ctx context.Context, messages []anthropic.MessageParam, profileJSON string) (<-chan string, <-chan error) {
	ch  := make(chan string, 64)
	ech := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(ech)

		sysPrompt := fmt.Sprintf(systemPromptCoach, profileJSON)
		stream, err := s.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
			Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
			MaxTokens: anthropic.F(int64(1024)),
			System: anthropic.F([]anthropic.TextBlockParam{
				{Text: anthropic.F(sysPrompt)},
			}),
			Messages: anthropic.F(messages),
		})
		if err != nil {
			ech <- err
			return
		}
		for stream.Next() {
			evt := stream.Current()
			switch v := evt.Delta.(type) {
			case anthropic.ContentBlockDeltaEventDelta:
				if v.Text != "" {
					ch <- v.Text
				}
			}
		}
		if err := stream.Err(); err != nil {
			ech <- err
		}
	}()

	return ch, ech
}

func buildAnalysisPrompt(req models.AnalysisRequest) string {
	skills := "not specified"
	if len(req.Skills) > 0 {
		skills = strings.Join(req.Skills, ", ")
	}
	loc := req.Location
	if loc == "" {
		loc = "not specified"
	}
	return fmt.Sprintf(`
Profession: %s
Current Skills: %s
Years of Experience: %d
Location: %s

Analyze the AI replacement risk for this professional and respond ONLY with valid JSON matching this exact schema:
{
  "profession": string,
  "ai_risk_score": integer 0-100,
  "risk_level": "Low"|"Medium"|"High"|"Critical",
  "explanation": string (2-3 sentences, honest and specific),
  "resilient_skills": [string] (3-5 skills AI cannot replace),
  "at_risk_tasks": [string] (3-5 tasks likely automated),
  "top_pivots": [
    {
      "title": string,
      "match_score": integer 0-100,
      "skill_gap": [string],
      "avg_reskill_months": integer,
      "reasoning": string,
      "avg_salary_usd": integer,
      "growth_pct": string (e.g. "+240%%")
    }
  ] (exactly 5 pivots),
  "timeline": string (when major disruption expected),
  "immediate_actions": [string] (exactly 3 actions for the next 6 months)
}
`, req.Profession, skills, req.YearsExp, loc)
}

const systemPromptAnalysis = `
You are FutureProof AI — the world's most advanced career intelligence system for the post-AI era.
Your analysis must be honest, data-backed, and specific. Never give vague platitudes.
Respond ONLY with valid JSON. No markdown. No explanation outside the JSON.
`

const systemPromptCoach = `
You are FutureProof AI Career Coach — empathetic, direct, and deeply knowledgeable about AI-driven career disruption.

User profile: %s

Guidelines:
- Give specific, actionable advice. Never vague platitudes.
- Be honest about automation risks without causing panic.
- Structure advice around: (1) immediate actions, (2) 6-month plan, (3) 2-year vision.
- Consider the user's location, experience, and skills in every response.
- Keep responses concise (under 300 words) unless depth is explicitly requested.
`

func slugify(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}

// ============================================================
// FILE: internal/services/career_service.go
// ============================================================

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/yourname/futureproof-api/internal/models"
	"github.com/yourname/futureproof-api/internal/repository"
	"github.com/yourname/futureproof-api/pkg/cache"
)

type CareerService struct {
	repo  *repository.CareerRepo
	cache *cache.Cache
}

func NewCareerService(repo *repository.CareerRepo, c *cache.Cache) *CareerService {
	return &CareerService{repo, c}
}

func (s *CareerService) GetAll(ctx context.Context) ([]models.Profession, error) {
	var out []models.Profession
	if err := s.cache.Get(ctx, "careers:all", &out); err == nil {
		return out, nil
	}
	out, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, "careers:all", out, 6*time.Hour)
	return out, nil
}

func (s *CareerService) GetBySlug(ctx context.Context, slug string) (*models.Profession, error) {
	key := fmt.Sprintf("career:%s", slug)
	var p models.Profession
	if err := s.cache.Get(ctx, key, &p); err == nil {
		return &p, nil
	}
	prof, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, key, prof, 6*time.Hour)
	return prof, nil
}

func (s *CareerService) GetTransitions(ctx context.Context, slug string) ([]models.CareerTransition, error) {
	key := fmt.Sprintf("transitions:%s", slug)
	var out []models.CareerTransition
	if err := s.cache.Get(ctx, key, &out); err == nil {
		return out, nil
	}
	out, err := s.repo.GetTransitions(ctx, slug)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, key, out, 12*time.Hour)
	return out, nil
}

func (s *CareerService) GetFutureJobs(ctx context.Context) ([]models.Profession, error) {
	var out []models.Profession
	if err := s.cache.Get(ctx, "careers:future", &out); err == nil {
		return out, nil
	}
	out, err := s.repo.GetFutureJobs(ctx)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, "careers:future", out, 6*time.Hour)
	return out, nil
}

// ============================================================
// FILE: internal/api/handlers/careers.go
// ============================================================

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourname/futureproof-api/internal/models"
	"github.com/yourname/futureproof-api/internal/services"
)

type CareerHandler struct {
	svc   *services.CareerService
	aiSvc *services.AIService
}

func NewCareerHandler(svc *services.CareerService, ai *services.AIService) *CareerHandler {
	return &CareerHandler{svc, ai}
}

func (h *CareerHandler) GetAll(c *fiber.Ctx) error {
	careers, err := h.svc.GetAll(c.Context())
	if err != nil {
		return fiber.NewError(500, "failed to fetch careers")
	}
	return c.JSON(fiber.Map{"data": careers, "total": len(careers)})
}

func (h *CareerHandler) GetBySlug(c *fiber.Ctx) error {
	prof, err := h.svc.GetBySlug(c.Context(), c.Params("slug"))
	if err != nil {
		return fiber.NewError(404, "profession not found")
	}
	return c.JSON(fiber.Map{"data": prof})
}

func (h *CareerHandler) GetFutureJobs(c *fiber.Ctx) error {
	jobs, err := h.svc.GetFutureJobs(c.Context())
	if err != nil {
		return fiber.NewError(500, "failed to fetch future jobs")
	}
	return c.JSON(fiber.Map{"data": jobs, "total": len(jobs)})
}

func (h *CareerHandler) GetTransitions(c *fiber.Ctx) error {
	transitions, err := h.svc.GetTransitions(c.Context(), c.Params("slug"))
	if err != nil {
		return fiber.NewError(404, "no transitions found")
	}
	return c.JSON(fiber.Map{"data": transitions})
}

func (h *CareerHandler) Analyze(c *fiber.Ctx) error {
	var req models.AnalysisRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid request body")
	}
	if req.Profession == "" {
		return fiber.NewError(400, "profession is required")
	}
	analysis, err := h.aiSvc.AnalyzeCareer(c.Context(), req)
	if err != nil {
		return fiber.NewError(500, "AI analysis failed: "+err.Error())
	}
	return c.JSON(fiber.Map{"data": analysis})
}

// ============================================================
// FILE: internal/api/handlers/cities.go
// ============================================================

package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/yourname/futureproof-api/internal/services"
)

type CityHandler struct{ svc *services.CityService }

func NewCityHandler(svc *services.CityService) *CityHandler { return &CityHandler{svc} }

func (h *CityHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "0"))
	cities, err := h.svc.GetAll(c.Context())
	if err != nil {
		return fiber.NewError(500, "failed to fetch cities")
	}
	if limit > 0 && limit < len(cities) {
		cities = cities[:limit]
	}
	return c.JSON(fiber.Map{"data": cities, "total": len(cities)})
}

func (h *CityHandler) GetByID(c *fiber.Ctx) error {
	city, err := h.svc.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return fiber.NewError(404, "city not found")
	}
	return c.JSON(fiber.Map{"data": city})
}

func (h *CityHandler) Search(c *fiber.Ctx) error {
	q := c.Query("q")
	if len(q) < 2 {
		return fiber.NewError(400, "query must be at least 2 characters")
	}
	cities, err := h.svc.Search(c.Context(), q)
	if err != nil {
		return fiber.NewError(500, "search failed")
	}
	return c.JSON(fiber.Map{"data": cities, "total": len(cities)})
}

// ============================================================
// FILE: internal/api/handlers/ai.go  (streaming chat)
// ============================================================

package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/gofiber/fiber/v2"
	"github.com/yourname/futureproof-api/internal/services"
)

type AIHandler struct{ aiSvc *services.AIService }

func NewAIHandler(ai *services.AIService) *AIHandler { return &AIHandler{ai} }

type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"     validate:"required"`
	UserProfile string        `json:"user_profile"` // JSON string of user context
}

type ChatMessage struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"`
}

func (h *AIHandler) Chat(c *fiber.Ctx) error {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid request body")
	}

	// Convert to Anthropic message params
	var msgs []anthropic.MessageParam
	for _, m := range req.Messages {
		if m.Role == "user" {
			msgs = append(msgs, anthropic.UserMessageParam(anthropic.NewTextBlock(m.Content)))
		} else {
			msgs = append(msgs, anthropic.AssistantMessageParam(anthropic.NewTextBlock(m.Content)))
		}
	}

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	textCh, errCh := h.aiSvc.StreamChat(c.Context(), msgs, req.UserProfile)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for {
			select {
			case chunk, ok := <-textCh:
				if !ok {
					fmt.Fprintf(w, "data: [DONE]\n\n")
					w.Flush()
					return
				}
				b, _ := json.Marshal(map[string]string{"delta": chunk})
				fmt.Fprintf(w, "data: %s\n\n", b)
				w.Flush()
			case err := <-errCh:
				if err != nil {
					b, _ := json.Marshal(map[string]string{"error": err.Error()})
					fmt.Fprintf(w, "data: %s\n\n", b)
					w.Flush()
				}
			}
		}
	})
	return nil
}

// ============================================================
// FILE: internal/api/router.go
// ============================================================

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourname/futureproof-api/internal/api/handlers"
	"github.com/yourname/futureproof-api/internal/api/middleware"
	"github.com/yourname/futureproof-api/internal/config"
	"github.com/yourname/futureproof-api/internal/repository"
	"github.com/yourname/futureproof-api/internal/services"
	"github.com/yourname/futureproof-api/pkg/cache"
	"github.com/yourname/futureproof-api/pkg/database"
	"go.uber.org/zap"
)

func SetupRoutes(app *fiber.App, db *database.DB, rdb *cache.Cache, cfg *config.Config, log *zap.Logger) {
	// Repos
	cityRepo   := repository.NewCityRepo(db)
	careerRepo := repository.NewCareerRepo(db)
	userRepo   := repository.NewUserRepo(db)

	// Services
	aiSvc     := services.NewAIService(cfg, rdb)
	citySvc   := services.NewCityService(cityRepo, rdb)
	careerSvc := services.NewCareerService(careerRepo, rdb)
	authSvc   := services.NewAuthService(userRepo, cfg)

	// Handlers
	cityH   := handlers.NewCityHandler(citySvc)
	careerH := handlers.NewCareerHandler(careerSvc, aiSvc)
	aiH     := handlers.NewAIHandler(aiSvc)
	authH   := handlers.NewAuthHandler(authSvc)

	// Rate limiter (Redis-backed)
	rl := middleware.NewRateLimiter(rdb)

	v1 := app.Group("/api/v1")

	// Health
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "version": "1.0.0"})
	})

	// Auth (public)
	auth := v1.Group("/auth")
	auth.Post("/register", authH.Register)
	auth.Post("/login",    authH.Login)
	auth.Post("/refresh",  authH.Refresh)

	// Cities (public, rate-limited)
	cities := v1.Group("/cities", rl.Limit(60))
	cities.Get("/",          cityH.GetAll)
	cities.Get("/search",    cityH.Search)
	cities.Get("/:id",       cityH.GetByID)

	// Careers (public, rate-limited)
	careers := v1.Group("/careers", rl.Limit(60))
	careers.Get("/",                  careerH.GetAll)
	careers.Get("/future",            careerH.GetFutureJobs)
	careers.Get("/:slug",             careerH.GetBySlug)
	careers.Get("/:slug/transitions", careerH.GetTransitions)

	// AI endpoints (rate-limited more strictly)
	ai := v1.Group("/ai", rl.Limit(10))
	ai.Post("/analyze", careerH.Analyze)
	ai.Post("/chat",    aiH.Chat)

	// Authenticated routes
	protected := v1.Group("/users", middleware.JWTAuth(cfg.JWTSecret))
	protected.Get("/me",            handlers.GetMe(userRepo))
	protected.Put("/me",            handlers.UpdateMe(userRepo))
}

// ErrorHandler is the global Fiber error handler
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg  := "internal server error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg  = e.Message
	}
	return c.Status(code).JSON(fiber.Map{"error": msg, "status": code})
}

// ============================================================
// FILE: internal/api/middleware/ratelimit.go
// ============================================================

package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yourname/futureproof-api/pkg/cache"
)

type RateLimiter struct{ cache *cache.Cache }

func NewRateLimiter(c *cache.Cache) *RateLimiter { return &RateLimiter{c} }

// Limit returns a middleware that allows `max` requests per minute per IP
func (rl *RateLimiter) Limit(max int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := fmt.Sprintf("rl:%s:%s", c.IP(), c.Path())
		ctx := context.Background()

		var count int
		rl.cache.Get(ctx, key, &count)
		if count >= max {
			return fiber.NewError(429, "rate limit exceeded — slow down")
		}
		count++
		_ = rl.cache.Set(ctx, key, count, time.Minute)
		return c.Next()
	}
}

// ============================================================
// FILE: internal/api/middleware/auth.go
// ============================================================

package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return fiber.NewError(401, "missing or invalid authorization header")
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(401, "unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(401, "invalid or expired token")
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		c.Locals("userID", claims["sub"])
		return c.Next()
	}
}

// ============================================================
// FILE: internal/workers/scheduler.go
// ============================================================

package workers

import (
	"github.com/robfig/cron/v3"
	"github.com/yourname/futureproof-api/internal/config"
	"github.com/yourname/futureproof-api/pkg/cache"
	"github.com/yourname/futureproof-api/pkg/database"
	"go.uber.org/zap"
)

type Scheduler struct {
	cron    *cron.Cron
	scraper *JobScraper
	log     *zap.Logger
}

func NewScheduler(db *database.DB, rdb *cache.Cache, cfg *config.Config, log *zap.Logger) *Scheduler {
	scraper := NewJobScraper(db, rdb, cfg, log)
	c := cron.New()

	s := &Scheduler{cron: c, scraper: scraper, log: log}

	// Refresh job market data every 6 hours
	c.AddFunc("0 */6 * * *", func() {
		log.Info("starting job market scrape")
		if err := scraper.Run(); err != nil {
			log.Error("job scrape failed", zap.Error(err))
		}
	})

	// Refresh World Bank data weekly
	c.AddFunc("0 2 * * 0", func() {
		log.Info("fetching World Bank data")
		if err := scraper.FetchWorldBankData(); err != nil {
			log.Error("world bank fetch failed", zap.Error(err))
		}
	})

	return s
}

func (s *Scheduler) Start() { s.cron.Start(); s.log.Info("scheduler started") }
func (s *Scheduler) Stop()  { s.cron.Stop() }

// ============================================================
// FILE: internal/workers/job_scraper.go
// ============================================================

package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/yourname/futureproof-api/internal/config"
	"github.com/yourname/futureproof-api/pkg/cache"
	"github.com/yourname/futureproof-api/pkg/database"
	"go.uber.org/zap"
)

type JobScraper struct {
	db     *database.DB
	cache  *cache.Cache
	cfg    *config.Config
	log    *zap.Logger
	client *http.Client
}

func NewJobScraper(db *database.DB, rdb *cache.Cache, cfg *config.Config, log *zap.Logger) *JobScraper {
	return &JobScraper{
		db: db, cache: rdb, cfg: cfg, log: log,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type AdzunaResponse struct {
	Results []struct {
		Title    string  `json:"title"`
		Location struct{ DisplayName string `json:"display_name"` } `json:"location"`
		SalaryMin float64 `json:"salary_min"`
		SalaryMax float64 `json:"salary_max"`
	} `json:"results"`
	Count int `json:"count"`
}

// Run triggers all data collection goroutines in parallel
func (s *JobScraper) Run() error {
	targets := []struct{ city, country string }{
		{"san-francisco", "us"}, {"new-york", "us"}, {"london", "gb"},
		{"berlin", "de"}, {"singapore", "sg"}, {"toronto", "ca"},
	}
	errs := make(chan error, len(targets))
	for _, t := range targets {
		go func(city, country string) {
			data, err := s.fetchAdzuna(city, country)
			if err != nil {
				errs <- err
				return
			}
			errs <- s.upsertJobData(context.Background(), city, data)
		}(t.city, t.country)
	}
	for range targets {
		if err := <-errs; err != nil {
			s.log.Warn("scrape target failed", zap.Error(err))
		}
	}
	return nil
}

func (s *JobScraper) fetchAdzuna(city, countryCode string) (*AdzunaResponse, error) {
	if s.cfg.AdzunaAppID == "" {
		return nil, fmt.Errorf("adzuna credentials not set")
	}
	apiURL := fmt.Sprintf(
		"https://api.adzuna.com/v1/api/jobs/%s/search/1?app_id=%s&app_key=%s&where=%s&results_per_page=50&sort_by=date",
		countryCode, s.cfg.AdzunaAppID, s.cfg.AdzunaKey, url.QueryEscape(city),
	)
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result AdzunaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *JobScraper) FetchWorldBankData() error {
	// World Bank API: employment ratios by country
	apiURL := "https://api.worldbank.org/v2/country/all/indicator/SL.EMP.TOTL.SP.ZS?format=json&per_page=50&mrv=1"
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	s.log.Info("world bank data fetched", zap.Int("status", resp.StatusCode))
	// TODO: parse and upsert into city_scores
	return nil
}

func (s *JobScraper) upsertJobData(ctx context.Context, city string, data *AdzunaResponse) error {
	s.log.Info("upserting job data", zap.String("city", city), zap.Int("count", data.Count))
	// Invalidate city caches after update
	_ = s.cache.Del(ctx, "cities:all", fmt.Sprintf("city:%s", city))
	return nil
}

// ============================================================
// FILE: .env.example
// ============================================================
//
// APP_ENV=development
// PORT=8080
// DATABASE_URL=postgres://fp:fp@localhost:5432/futureproof?sslmode=disable
// REDIS_URL=redis://localhost:6379
// ANTHROPIC_API_KEY=sk-ant-...
// JWT_SECRET=your-super-secret-jwt-key-change-in-production
// ALLOWED_ORIGINS=http://localhost:3000
// ADZUNA_APP_ID=
// ADZUNA_KEY=

// ============================================================
// FILE: Dockerfile
// ============================================================
//
// FROM golang:1.22-alpine AS builder
// WORKDIR /app
// COPY go.mod go.sum ./
// RUN go mod download
// COPY . .
// RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/server
//
// FROM scratch
// COPY --from=builder /app/server /server
// COPY --from=builder /app/.env.example /.env.example
// EXPOSE 8080
// ENTRYPOINT ["/server"]

// ============================================================
// FILE: docker-compose.yml
// ============================================================
//
// version: "3.9"
// services:
//   api:
//     build: .
//     ports: ["8080:8080"]
//     env_file: .env
//     depends_on: [postgres, redis]
//     restart: unless-stopped
//
//   postgres:
//     image: postgres:16-alpine
//     environment:
//       POSTGRES_DB: futureproof
//       POSTGRES_USER: fp
//       POSTGRES_PASSWORD: fp
//     volumes: [pgdata:/var/lib/postgresql/data]
//     ports: ["5432:5432"]
//
//   redis:
//     image: redis:7-alpine
//     volumes: [redisdata:/data]
//     ports: ["6379:6379"]
//
// volumes:
//   pgdata:
//   redisdata:
