package handler

import (
	"net/http"

	"github.com/andrimuhayat/crud-test/internal/domain"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related endpoints.
type AuthHandler struct {
	authUC domain.AuthUseCase
}

// NewAuthHandler wires the handler to the auth use-case.
func NewAuthHandler(authUC domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

type tokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GenerateToken handles POST /auth/token.
// Level 5 requirement: accept credentials, return a signed JWT.
func (h *AuthHandler) GenerateToken(c *fiber.Ctx) error {
	var req tokenRequest
	if err := c.BodyParser(&req); err != nil || req.Username == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "username and password are required"})
	}

	token, err := h.authUC.GenerateToken(req.Username, req.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return c.JSON(fiber.Map{"token": token})
}
