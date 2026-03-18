package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Anthropic  AnthropicConfig
	ExternalAPI ExternalAPIConfig
	CORS       CORSConfig
	RateLimit  RateLimitConfig
	WebSocket  WebSocketConfig
	Workers    WorkersConfig
}

type AppConfig struct {
	Name string
	Port int
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	MaxConns int
	MinConns int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret        string
	Expiry        time.Duration
	RefreshExpiry time.Duration
}

type AnthropicConfig struct {
	APIKey string
	Model  string
}

type ExternalAPIConfig struct {
	AdzunaAppID     string
	AdzunaAPIKey    string
	WorldBankAPIURL string
	BLSAPIKey       string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type RateLimitConfig struct {
	Requests int
	Window   int
}

type WebSocketConfig struct {
	PingInterval time.Duration
	PongTimeout  time.Duration
}

type WorkersConfig struct {
	ScraperInterval   time.Duration
	DataFetchInterval time.Duration
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "FutureProof"),
			Port: getEnvInt("PORT", 8080),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "futureproof"),
			Password: getEnv("DB_PASSWORD", "futureproof_dev_pass"),
			Name:     getEnv("DB_NAME", "futureproof_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: getEnvInt("DB_MAX_CONNS", 25),
			MinConns: getEnvInt("DB_MIN_CONNS", 5),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "change-this-secret-key"),
			Expiry:        getEnvDuration("JWT_EXPIRY", 24*time.Hour),
			RefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 168*time.Hour),
		},
		Anthropic: AnthropicConfig{
			APIKey: getEnv("ANTHROPIC_API_KEY", ""),
			Model:  getEnv("ANTHROPIC_MODEL", "claude-3-5-sonnet-20241022"),
		},
		ExternalAPI: ExternalAPIConfig{
			AdzunaAppID:     getEnv("ADZUNA_APP_ID", ""),
			AdzunaAPIKey:    getEnv("ADZUNA_API_KEY", ""),
			WorldBankAPIURL: getEnv("WORLD_BANK_API_URL", "https://api.worldbank.org/v2"),
			BLSAPIKey:       getEnv("BLS_API_KEY", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:5173", "http://localhost:3000"}),
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getEnvInt("RATE_LIMIT_WINDOW", 60),
		},
		WebSocket: WebSocketConfig{
			PingInterval: getEnvDuration("WS_PING_INTERVAL", 30*time.Second),
			PongTimeout:  getEnvDuration("WS_PONG_TIMEOUT", 60*time.Second),
		},
		Workers: WorkersConfig{
			ScraperInterval:   getEnvDuration("SCRAPER_INTERVAL", 6*time.Hour),
			DataFetchInterval: getEnvDuration("DATA_FETCH_INTERVAL", 24*time.Hour),
		},
	}

	// Validate critical configuration
	if cfg.Anthropic.APIKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is required")
	}

	return cfg, nil
}

// Helper functions to read environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma
		result := []string{}
		current := ""
		for _, char := range value {
			if char == ',' {
				if current != "" {
					result = append(result, current)
					current = ""
				}
			} else {
				current += string(char)
			}
		}
		if current != "" {
			result = append(result, current)
		}
		return result
	}
	return defaultValue
}
