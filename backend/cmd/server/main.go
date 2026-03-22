package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/api"
	"github.com/shaykhnazar/futureproof/internal/api/middleware"
	"github.com/shaykhnazar/futureproof/internal/config"
	"github.com/shaykhnazar/futureproof/internal/repository"
	"github.com/shaykhnazar/futureproof/internal/services"
	"github.com/shaykhnazar/futureproof/pkg/cache"
	"github.com/shaykhnazar/futureproof/pkg/database"
	"github.com/shaykhnazar/futureproof/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	// Initialize logger
	zapLogger, err := logger.NewLogger(os.Getenv("ENV"))
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		zapLogger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize database
	zapLogger.Info("Connecting to database...",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Name),
	)
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Redis
	zapLogger.Info("Connecting to Redis...",
		zap.String("host", cfg.Redis.Host),
		zap.Int("port", cfg.Redis.Port),
	)
	redisClient := cache.NewRedis(cfg.Redis)
	defer redisClient.Close()

	// Ping Redis to verify connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx); err != nil {
		zapLogger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	// Initialize repositories
	careerRepo := repository.NewCareerRepository(db)
	cityRepo := repository.NewCityRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	careerService := services.NewCareerService(careerRepo, redisClient, zapLogger)
	cityService := services.NewCityService(cityRepo, redisClient, zapLogger)
	aiService := services.NewAIService(cfg.Anthropic.APIKey, cfg.Anthropic.Model, redisClient, zapLogger)
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.Expiry, zapLogger)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: middleware.ErrorHandler(zapLogger),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.Logger(zapLogger))
	app.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// Setup API routes
	api.SetupRoutes(app, api.RouterConfig{
		CareerService: careerService,
		CityService:   cityService,
		AIService:     aiService,
		AuthService:   authService,
		JWTSecret:     cfg.JWT.Secret,
		RateLimit: api.RateLimitConfig{
			Requests: cfg.RateLimit.Requests,
			Window:   cfg.RateLimit.Window,
		},
		Logger: zapLogger,
		Redis:  redisClient,
	})

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		zapLogger.Info("Server starting",
			zap.String("address", addr),
			zap.String("env", cfg.App.Env),
		)
		if err := app.Listen(addr); err != nil {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		zapLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	zapLogger.Info("Server stopped gracefully")
}
