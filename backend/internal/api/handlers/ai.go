package handlers

import (
	"bufio"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/internal/services"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	service *services.AIService
	logger  *zap.Logger
}

// NewAIHandler creates a new AI handler
func NewAIHandler(service *services.AIService, logger *zap.Logger) *AIHandler {
	return &AIHandler{
		service: service,
		logger:  logger,
	}
}

// AnalyzeCareer handles POST /api/v1/analyze
func (h *AIHandler) AnalyzeCareer(c *fiber.Ctx) error {
	var req models.AnalysisRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.ProfessionSlug == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Profession slug is required")
	}

	// Perform AI analysis
	result, err := h.service.AnalyzeCareer(c.Context(), req)
	if err != nil {
		h.logger.Error("Failed to analyze career", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to analyze career")
	}

	return c.JSON(result)
}

// ChatWithCoach handles POST /api/v1/ai/chat (streaming)
func (h *AIHandler) ChatWithCoach(c *fiber.Ctx) error {
	var req struct {
		Message string `json:"message"`
		History []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"history"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Message == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Message is required")
	}

	// Set headers for Server-Sent Events
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Convert history to role/content pairs
	history := make([]services.ChatMessage, len(req.History))
	for i, h := range req.History {
		history[i] = services.ChatMessage{Role: h.Role, Content: h.Content}
	}

	// Stream response
	responseChan, err := h.service.StreamChat(c.Context(), req.Message, history)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to start chat stream")
	}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for chunk := range responseChan {
			_, _ = w.WriteString("data: " + chunk + "\n\n")
			_ = w.Flush()
		}
		_, _ = w.WriteString("data: [DONE]\n\n")
		_ = w.Flush()
	})

	return nil
}
