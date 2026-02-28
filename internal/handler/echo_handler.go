package handler

import "github.com/gofiber/fiber/v2"

// EchoHandler echoes the request JSON body back to the caller.
type EchoHandler struct{}

// NewEchoHandler returns an EchoHandler.
func NewEchoHandler() *EchoHandler { return &EchoHandler{} }

// Echo handles POST /echo.
// Level 2 requirement: reflect the raw request body byte-for-byte as the
// response body, preserving key order, whitespace, and all formatting exactly
// as the client sent them.
//
// The body is never deserialised into a map; doing so would cause encoding/json
// to re-serialise map keys in alphabetical order, breaking strict string
// comparisons on the test platform.
func (h *EchoHandler) Echo(c *fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "empty body"})
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Send(body)
}
