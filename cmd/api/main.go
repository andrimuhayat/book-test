// Package main is the entry point for the API Quest backend service.
package main

import (
	"log"

	"github.com/andrimuhayat/crud-test/internal/handler"
	"github.com/andrimuhayat/crud-test/internal/repository/memory"
	"github.com/andrimuhayat/crud-test/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// --- Dependency wiring (composition root) ---
	bookRepo := memory.NewBookRepository()
	bookUC := usecase.NewBookUseCase(bookRepo)
	authUC := usecase.NewAuthUseCase()

	pingH := handler.NewPingHandler()
	echoH := handler.NewEchoHandler()
	authH := handler.NewAuthHandler(authUC)
	bookH := handler.NewBookHandler(bookUC)

	// --- Fiber app ---
	app := fiber.New(fiber.Config{
		// Return errors as JSON instead of plain text.
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	// --- Public routes ---
	app.Get("/ping", pingH.Ping)
	app.Post("/echo", echoH.Echo)
	app.Post("/auth/token", authH.GenerateToken)

	// --- Book routes (public — no auth required for Level 3) ---
	app.Post("/books", bookH.CreateBook)
	app.Get("/books", bookH.GetBooks)
	app.Get("/books/:id", bookH.GetBook)
	app.Put("/books/:id", bookH.UpdateBook)
	app.Delete("/books/:id", bookH.DeleteBook)

	log.Fatal(app.Listen(":8080"))
}
