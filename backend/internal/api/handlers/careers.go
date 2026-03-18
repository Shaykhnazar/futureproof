package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/api/middleware"
	"github.com/shaykhnazar/futureproof/internal/services"
)

// CareerHandler handles career-related HTTP requests
type CareerHandler struct {
	service *services.CareerService
	logger  *zap.Logger
}

// NewCareerHandler creates a new career handler
func NewCareerHandler(service *services.CareerService, logger *zap.Logger) *CareerHandler {
	return &CareerHandler{
		service: service,
		logger:  logger,
	}
}

// GetAllProfessions handles GET /api/v1/professions
func (h *CareerHandler) GetAllProfessions(c *fiber.Ctx) error {
	professions, err := h.service.GetAllProfessions(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve professions")
	}

	return c.JSON(fiber.Map{
		"professions": professions,
		"count":       len(professions),
	})
}

// GetProfessionBySlug handles GET /api/v1/professions/:slug
func (h *CareerHandler) GetProfessionBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Profession slug is required")
	}

	profession, err := h.service.GetProfessionBySlug(c.Context(), slug)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Profession not found")
	}

	return c.JSON(profession)
}

// GetCareerTransitions handles GET /api/v1/professions/:slug/pivots
func (h *CareerHandler) GetCareerTransitions(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Profession slug is required")
	}

	transitions, err := h.service.GetCareerTransitions(c.Context(), slug)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve career transitions")
	}

	return c.JSON(fiber.Map{
		"transitions": transitions,
		"count":       len(transitions),
	})
}

// GetFutureProfessions handles GET /api/v1/professions/future
func (h *CareerHandler) GetFutureProfessions(c *fiber.Ctx) error {
	professions, err := h.service.GetFutureProfessions(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve future professions")
	}

	return c.JSON(fiber.Map{
		"professions": professions,
		"count":       len(professions),
	})
}

// SaveCareer handles POST /api/v1/careers/save (protected)
func (h *CareerHandler) SaveCareer(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var req struct {
		ProfessionID string `json:"profession_id"`
		Notes        string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	professionID, err := uuid.Parse(req.ProfessionID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid profession ID")
	}

	err = h.service.SaveCareer(c.Context(), userID, professionID, req.Notes)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save career")
	}

	return c.JSON(fiber.Map{
		"message": "Career saved successfully",
	})
}

// GetSavedCareers handles GET /api/v1/careers/saved (protected)
func (h *CareerHandler) GetSavedCareers(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	professions, err := h.service.GetSavedCareers(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve saved careers")
	}

	return c.JSON(fiber.Map{
		"professions": professions,
		"count":       len(professions),
	})
}
