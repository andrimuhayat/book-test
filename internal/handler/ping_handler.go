// Package handler contains the HTTP presentation layer (adapters).
package handler

import "github.com/gofiber/fiber/v2"

// PingHandler serves the health-check endpoint.
type PingHandler struct{}

// NewPingHandler returns a PingHandler.
func NewPingHandler() *PingHandler { return &PingHandler{} }

// Ping handles GET /ping.
// Level 1 requirement: returns {"message":"pong"} with HTTP 200.
func (h *PingHandler) Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "pong"})
}
