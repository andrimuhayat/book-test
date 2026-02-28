// Package memory provides thread-safe in-memory implementations of domain repositories.
package memory

import (
	"sync"

	"github.com/andrimuhayat/crud-test/internal/domain"
)

// BookRepository is a thread-safe, in-memory implementation of domain.BookRepository.
// It uses a sync.RWMutex to allow many concurrent readers but only one writer at a time,
// which is efficient for read-heavy workloads.
type BookRepository struct {
	mu    sync.RWMutex
	books map[string]*domain.Book
	order []string // insertion-order slice of IDs for stable LIST results
}

// NewBookRepository creates and returns an initialised BookRepository.
func NewBookRepository() *BookRepository {
	return &BookRepository{
		books: make(map[string]*domain.Book),
		order: make([]string, 0),
	}
}

// Create stores a new book. O(1) amortised.
func (r *BookRepository) Create(book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[book.ID] = book
	r.order = append(r.order, book.ID)
	return nil
}

// GetByID returns a single book by ID. Returns domain.ErrNotFound if absent.
func (r *BookRepository) GetByID(id string) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.books[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return book, nil
}

// GetAll returns books matching the filter, plus the total count before pagination.
// Filtering by Author is case-sensitive substring-free (exact match).
// Pagination uses 1-based page numbers.
func (r *BookRepository) GetAll(filter domain.BookFilter) ([]*domain.Book, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*domain.Book, 0, len(r.order))
	for _, id := range r.order {
		book := r.books[id]
		if filter.Author != "" && book.Author != filter.Author {
			continue
		}
		filtered = append(filtered, book)
	}

	total := len(filtered)

	// Apply pagination only when both page and limit are explicitly provided.
	if filter.Page > 0 && filter.Limit > 0 {
		start := (filter.Page - 1) * filter.Limit
		if start >= total {
			return make([]*domain.Book, 0), total, nil
		}
		end := start + filter.Limit
		if end > total {
			end = total
		}
		filtered = filtered[start:end]
	}

	return filtered, total, nil
}

// Update replaces the stored book. Returns domain.ErrNotFound if the ID is absent.
func (r *BookRepository) Update(book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[book.ID]; !ok {
		return domain.ErrNotFound
	}
	r.books[book.ID] = book
	return nil
}

// Delete removes a book by ID. Returns domain.ErrNotFound if the ID is absent.
func (r *BookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.books, id)

	// Remove id from the insertion-order slice.
	for i, oid := range r.order {
		if oid == id {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	return nil
}

