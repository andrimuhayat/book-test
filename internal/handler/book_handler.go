package handler

import (
	"net/http"

	"github.com/andrimuhayat/crud-test/internal/domain"
	"github.com/gofiber/fiber/v2"
)

// BookHandler handles CRUD and search endpoints for books.
type BookHandler struct {
	bookUC domain.BookUseCase
}

// NewBookHandler wires the handler to the book use-case.
func NewBookHandler(bookUC domain.BookUseCase) *BookHandler {
	return &BookHandler{bookUC: bookUC}
}

type createBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// CreateBook handles POST /books.
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var req createBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Title == "" || req.Author == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "title and author are required"})
	}

	book, err := h.bookUC.CreateBook(req.Title, req.Author, req.Year)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(book)
}

// GetBook handles GET /books/:id.
func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	book, err := h.bookUC.GetBook(id)
	if err == domain.ErrNotFound {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(book)
}

// GetBooks handles GET /books with optional ?author= query param.
// Returns a bare JSON array of book objects (Level 3 requirement).
func (h *BookHandler) GetBooks(c *fiber.Ctx) error {
	filter := domain.BookFilter{
		Author: c.Query("author"),
		Page:   1,
		Limit:  1000,
	}

	books, _, err := h.bookUC.GetBooks(filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if books == nil {
		books = []*domain.Book{}
	}
	return c.JSON(books)
}

// UpdateBook handles PUT /books/:id.
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var req createBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Title == "" || req.Author == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "title and author are required"})
	}

	book, err := h.bookUC.UpdateBook(id, req.Title, req.Author, req.Year)
	if err == domain.ErrNotFound {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
	}
	if err == domain.ErrInvalidData {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid data"})
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(book)
}

// DeleteBook handles DELETE /books/:id.
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.bookUC.DeleteBook(id)
	if err == domain.ErrNotFound {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(http.StatusNoContent)
}
