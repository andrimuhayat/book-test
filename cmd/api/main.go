// Package main is the entry point for the API Quest backend service.
package main

import (
	"log"

	"github.com/andrimuhayat/crud-test/internal/handler"
	"github.com/andrimuhayat/crud-test/internal/middleware"
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

	// --- Protected book routes (Level 5 — JWT required) ---
	books := app.Group("/books", middleware.Auth(authUC))
	books.Post("/", bookH.CreateBook)
	books.Get("/", bookH.GetBooks)
	books.Get("/:id", bookH.GetBook)
	books.Put("/:id", bookH.UpdateBook)
	books.Delete("/:id", bookH.DeleteBook)

	log.Fatal(app.Listen(":8080"))
}
