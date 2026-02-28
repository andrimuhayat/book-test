package domain

import "time"

// Book represents the core book entity.
type Book struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int       `json:"year,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// BookFilter holds query parameters for listing books.
type BookFilter struct {
	Author string
	Page   int
	Limit  int
}

// BookRepository defines the persistence contract for books.
// Implementations must be safe for concurrent use.
type BookRepository interface {
	Create(book *Book) error
	GetByID(id string) (*Book, error)
	GetAll(filter BookFilter) ([]*Book, int, error)
	Update(book *Book) error
	Delete(id string) error
}

// BookUseCase defines the business-logic contract for books.
type BookUseCase interface {
	CreateBook(title, author string, year int) (*Book, error)
	GetBook(id string) (*Book, error)
	GetBooks(filter BookFilter) ([]*Book, int, error)
	UpdateBook(id, title, author string, year int) (*Book, error)
	DeleteBook(id string) error
}

