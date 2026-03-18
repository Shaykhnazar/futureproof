package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/api/handlers"
	"github.com/shaykhnazar/futureproof/internal/api/middleware"
	"github.com/shaykhnazar/futureproof/internal/services"
	"github.com/shaykhnazar/futureproof/pkg/cache"
)

// RouterConfig holds dependencies for routing setup
type RouterConfig struct {
	CareerService *services.CareerService
	CityService   *services.CityService
	AIService     *services.AIService
	AuthService   *services.AuthService
	JWTSecret     string
	RateLimit     RateLimitConfig
	Logger        *zap.Logger
	Redis         *cache.Redis
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int
	Window   int
}

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, cfg RouterConfig) {
	// Initialize handlers
	careerHandler := handlers.NewCareerHandler(cfg.CareerService, cfg.Logger)
	cityHandler := handlers.NewCityHandler(cfg.CityService, cfg.Logger)
	aiHandler := handlers.NewAIHandler(cfg.AIService, cfg.Logger)
	userHandler := handlers.NewUserHandler(cfg.AuthService, cfg.Logger)
	wsHub := handlers.NewWebSocketHub(cfg.Logger)

	// API v1 group
	api := app.Group("/api/v1")

	// Apply rate limiting to API routes
	api.Use(middleware.RateLimit(cfg.Redis, cfg.RateLimit.Requests, cfg.RateLimit.Window))

	// Public routes

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// Cities routes
	cities := api.Group("/cities")
	cities.Get("/", cityHandler.GetAllCities)
	cities.Get("/:id", cityHandler.GetCityByID)
	cities.Get("/region/:region", cityHandler.GetCitiesByRegion)

	// Professions routes
	professions := api.Group("/professions")
	professions.Get("/", careerHandler.GetAllProfessions)
	professions.Get("/future", careerHandler.GetFutureProfessions)
	professions.Get("/:slug", careerHandler.GetProfessionBySlug)
	professions.Get("/:slug/pivots", careerHandler.GetCareerTransitions)

	// AI Analysis routes
	api.Post("/analyze", aiHandler.AnalyzeCareer)
	api.Post("/ai/chat", aiHandler.ChatWithCoach)

	// Protected routes (require authentication)
	authMiddleware := middleware.Auth(cfg.AuthService)

	// User routes (protected)
	users := api.Group("/users", authMiddleware)
	users.Get("/me", userHandler.GetCurrentUser)
	users.Put("/profile", userHandler.UpdateProfile)

	// Career management routes (protected)
	careers := api.Group("/careers", authMiddleware)
	careers.Post("/save", careerHandler.SaveCareer)
	careers.Get("/saved", careerHandler.GetSavedCareers)

	// WebSocket route for real-time updates
	app.Get("/ws/globe", handlers.UpgradeMiddleware, websocket.New(wsHub.HandleWebSocket))

	// Example: Trigger a globe update (for testing)
	api.Post("/trigger-update", func(c *fiber.Ctx) error {
		wsHub.Broadcast("city_update", fiber.Map{
			"city_name": "San Francisco",
			"new_score": 96,
			"reason":    "New AI investment announced",
		})
		return c.JSON(fiber.Map{
			"message": "Update broadcast sent",
		})
	})

	cfg.Logger.Info("Routes configured successfully")
}
