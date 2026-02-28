// Package middleware provides Fiber middleware functions.
package middleware

import (
	"strings"

	"github.com/andrimuhayat/crud-test/internal/domain"
	"github.com/gofiber/fiber/v2"
)

const usernameLocalKey = "username"

// Auth returns a Fiber middleware that validates a Bearer JWT.
// On success the parsed subject claim is stored in c.Locals("username").
func Auth(authUC domain.AuthUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
		}

		username, err := authUC.ValidateToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		c.Locals(usernameLocalKey, username)
		return c.Next()
	}
}

