// Package usecase implements the application business logic layer.
package usecase

import (
	"time"

	"github.com/andrimuhayat/crud-test/internal/domain"
	"github.com/google/uuid"
)

// BookUseCase implements domain.BookUseCase.
type BookUseCase struct {
	repo domain.BookRepository
}

// NewBookUseCase wires the use-case to a repository.
func NewBookUseCase(repo domain.BookRepository) *BookUseCase {
	return &BookUseCase{repo: repo}
}

// CreateBook validates input, assigns a UUID, and persists a new book.
func (uc *BookUseCase) CreateBook(title, author string, year int) (*domain.Book, error) {
	if title == "" || author == "" {
		return nil, domain.ErrInvalidData
	}

	book := &domain.Book{
		ID:        uuid.New().String(),
		Title:     title,
		Author:    author,
		Year:      year,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.repo.Create(book); err != nil {
		return nil, err
	}
	return book, nil
}

// GetBook retrieves a book by ID.
func (uc *BookUseCase) GetBook(id string) (*domain.Book, error) {
	return uc.repo.GetByID(id)
}

// GetBooks returns a (optionally filtered, optionally paginated) list of books
// together with the total count before pagination.
func (uc *BookUseCase) GetBooks(filter domain.BookFilter) ([]*domain.Book, int, error) {
	return uc.repo.GetAll(filter)
}

// UpdateBook replaces the mutable fields of an existing book.
func (uc *BookUseCase) UpdateBook(id, title, author string, year int) (*domain.Book, error) {
	if title == "" || author == "" {
		return nil, domain.ErrInvalidData
	}

	existing, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	existing.Title = title
	existing.Author = author
	existing.Year = year

	if err := uc.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteBook removes a book by ID.
func (uc *BookUseCase) DeleteBook(id string) error {
	return uc.repo.Delete(id)
}

