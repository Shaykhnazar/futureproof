package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/services"
)

// CityHandler handles city-related HTTP requests
type CityHandler struct {
	service *services.CityService
	logger  *zap.Logger
}

// NewCityHandler creates a new city handler
func NewCityHandler(service *services.CityService, logger *zap.Logger) *CityHandler {
	return &CityHandler{
		service: service,
		logger:  logger,
	}
}

// GetAllCities handles GET /api/v1/cities
func (h *CityHandler) GetAllCities(c *fiber.Ctx) error {
	// Optional region filter
	region := c.Query("region")

	var cities interface{}
	var err error

	if region != "" {
		cities, err = h.service.GetCitiesByRegion(c.Context(), region)
	} else {
		cities, err = h.service.GetAllCities(c.Context())
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve cities")
	}

	return c.JSON(fiber.Map{
		"cities": cities,
	})
}

// GetCityByID handles GET /api/v1/cities/:id
func (h *CityHandler) GetCityByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "City ID is required")
	}

	cityID, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid city ID")
	}

	city, err := h.service.GetCityByID(c.Context(), cityID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "City not found")
	}

	return c.JSON(city)
}

// GetCitiesByRegion handles GET /api/v1/cities/region/:region
func (h *CityHandler) GetCitiesByRegion(c *fiber.Ctx) error {
	region := c.Params("region")
	if region == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Region is required")
	}

	cities, err := h.service.GetCitiesByRegion(c.Context(), region)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve cities")
	}

	return c.JSON(fiber.Map{
		"cities": cities,
		"region": region,
		"count":  len(cities),
	})
}
