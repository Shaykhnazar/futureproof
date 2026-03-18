package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/api/middleware"
	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/internal/services"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service *services.AuthService
	logger  *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *services.AuthService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// Register handles POST /api/v1/auth/register
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	if len(req.Password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 8 characters")
	}

	// Register user
	user, err := h.service.Register(c.Context(), req)
	if err != nil {
		h.logger.Error("Registration failed", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login handles POST /api/v1/auth/login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	// Login user
	response, err := h.service.Login(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	return c.JSON(response)
}

// GetCurrentUser handles GET /api/v1/users/me (protected)
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	user, err := h.service.GetUserWithProfile(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.JSON(user)
}

// UpdateProfile handles PUT /api/v1/users/profile (protected)
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var profile models.UserProfile
	if err := c.BodyParser(&profile); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Set user ID from context
	profile.UserID = userID

	// Update profile
	err = h.service.UpdateProfile(c.Context(), &profile)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update profile")
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}
