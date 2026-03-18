package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorHandler creates a global error handler
func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default error code
		code := fiber.StatusInternalServerError
		message := "Internal server error"

		// Check if it's a Fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		// Log the error
		logger.Error("Request error",
			zap.Int("status", code),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Error(err),
		)

		// Send error response
		return c.Status(code).JSON(fiber.Map{
			"error": message,
			"code":  code,
		})
	}
}
