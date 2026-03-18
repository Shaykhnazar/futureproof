package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shaykhnazar/futureproof/pkg/cache"
)

// RateLimit creates a rate limiting middleware using Redis
func RateLimit(redis *cache.Redis, requests int, window int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client IP
		ip := c.IP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		// Increment counter
		count, err := redis.Incr(c.Context(), key)
		if err != nil {
			// If Redis fails, allow the request
			return c.Next()
		}

		// Set expiration on first request
		if count == 1 {
			_ = redis.Expire(c.Context(), key, time.Duration(window)*time.Second)
		}

		// Check if rate limit exceeded
		if count > int64(requests) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"retry_after": window,
			})
		}

		// Add rate limit headers
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", requests-int(count)))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Duration(window)*time.Second).Unix()))

		return c.Next()
	}
}
