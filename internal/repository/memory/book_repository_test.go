package memory_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/andrimuhayat/crud-test/internal/domain"
	"github.com/andrimuhayat/crud-test/internal/repository/memory"
)

func newBook(i int) *domain.Book {
	return &domain.Book{
		ID:        fmt.Sprintf("book-%d", i),
		Title:     fmt.Sprintf("Title %d", i),
		Author:    fmt.Sprintf("Author %d", i),
		Year:      2000 + i,
		CreatedAt: time.Now().UTC(),
	}
}

// TestConcurrentWrites verifies that many goroutines can Create books concurrently
// without data loss or races (run with -race).
func TestConcurrentWrites(t *testing.T) {
	const n = 200
	repo := memory.NewBookRepository()
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := repo.Create(newBook(i)); err != nil {
				t.Errorf("Create(%d) unexpected error: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	books, total, err := repo.GetAll(domain.BookFilter{})
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if total != n {
		t.Errorf("expected %d books, got %d", n, total)
	}
	if len(books) != n {
		t.Errorf("expected slice length %d, got %d", n, len(books))
	}
}

// TestConcurrentReads verifies that concurrent readers never block each other.
func TestConcurrentReads(t *testing.T) {
	const n = 50
	repo := memory.NewBookRepository()
	for i := 0; i < n; i++ {
		_ = repo.Create(newBook(i))
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("book-%d", i)
			if _, err := repo.GetByID(id); err != nil {
				t.Errorf("GetByID(%s) unexpected error: %v", id, err)
			}
		}(i)
	}
	wg.Wait()
}

// TestConcurrentMixedOps runs concurrent creates, reads, updates and deletes to
// stress-test the mutex strategy.
func TestConcurrentMixedOps(t *testing.T) {
	const n = 100
	repo := memory.NewBookRepository()

	// Pre-populate
	for i := 0; i < n; i++ {
		_ = repo.Create(newBook(i))
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		i := i
		// reader
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, _ = repo.GetAll(domain.BookFilter{Page: 1, Limit: 10})
		}()

		// updater
		wg.Add(1)
		go func() {
			defer wg.Done()
			b := newBook(i)
			b.Title = "Updated " + b.Title
			_ = repo.Update(b)
		}()

		// deleter (odd indices only to leave some books)
		if i%2 == 1 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = repo.Delete(fmt.Sprintf("book-%d", i))
			}()
		}
	}
	wg.Wait()
}

// TestPagination verifies correct page/total behaviour.
func TestPagination(t *testing.T) {
	repo := memory.NewBookRepository()
	for i := 0; i < 25; i++ {
		_ = repo.Create(newBook(i))
	}

	page1, total, err := repo.GetAll(domain.BookFilter{Page: 1, Limit: 10})
	if err != nil || total != 25 || len(page1) != 10 {
		t.Errorf("page1: got len=%d total=%d err=%v, want len=10 total=25 err=nil", len(page1), total, err)
	}

	page3, _, _ := repo.GetAll(domain.BookFilter{Page: 3, Limit: 10})
	if len(page3) != 5 {
		t.Errorf("page3: got len=%d, want 5", len(page3))
	}

	beyond, _, _ := repo.GetAll(domain.BookFilter{Page: 10, Limit: 10})
	if len(beyond) != 0 {
		t.Errorf("beyond: got len=%d, want 0", len(beyond))
	}
}

// TestNotFound ensures domain.ErrNotFound is returned for missing IDs.
func TestNotFound(t *testing.T) {
	repo := memory.NewBookRepository()
	if _, err := repo.GetByID("missing"); err != domain.ErrNotFound {
		t.Errorf("GetByID: want ErrNotFound, got %v", err)
	}
	if err := repo.Update(&domain.Book{ID: "missing"}); err != domain.ErrNotFound {
		t.Errorf("Update: want ErrNotFound, got %v", err)
	}
	if err := repo.Delete("missing"); err != domain.ErrNotFound {
		t.Errorf("Delete: want ErrNotFound, got %v", err)
	}
}
