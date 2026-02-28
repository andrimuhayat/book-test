package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// EchoHandler echoes the request JSON body back to the caller.
type EchoHandler struct{}

// NewEchoHandler returns an EchoHandler.
func NewEchoHandler() *EchoHandler { return &EchoHandler{} }

// Echo handles POST /echo.
// Level 2 requirement: reflect the full request body as the response body.
func (h *EchoHandler) Echo(c *fiber.Ctx) error {
	var body map[string]any
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}
	return c.JSON(body)
}
