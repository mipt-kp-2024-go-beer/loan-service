package repo

import (
	"context"
	"maps"
	"sync"
	"time"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
)

func NewMemoryRepo(dsn string) TestMemoryRepo {
	return &memoryRepo{
		mutex:     sync.RWMutex{},
		lentBooks: make(map[string]loans.LentBook),
	}
}

type memoryRepo struct {
	mutex     sync.RWMutex
	lentBooks map[string]loans.LentBook
}

// TestMemoryRepo is an interface that exposes memoryRepo's internal methods
// for use in unit testing
type TestMemoryRepo interface {
	loans.Repo
	LookupBook(ctx context.Context, ID string) (loans.LentBook, error)
	UpdateBook(ctx context.Context, book loans.LentBook) error
	InsertBook(ctx context.Context, book loans.LentBook) error
	RawData() map[string]loans.LentBook
}

func (m *memoryRepo) LookupBook(ctx context.Context, ID string) (loans.LentBook, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	book, ok := m.lentBooks[ID]
	if !ok {
		return loans.LentBook{}, fail.ErrNotFound
	}
	return book, nil
}

func (m *memoryRepo) UpdateBook(ctx context.Context, book loans.LentBook) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.lentBooks[book.ID]; !ok {
		return fail.ErrNotFound
	}
	m.lentBooks[book.ID] = book
	return nil
}

func (m *memoryRepo) InsertBook(ctx context.Context, book loans.LentBook) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.lentBooks[book.ID]; ok {
		return fail.ErrCollision
	}
	m.lentBooks[book.ID] = book
	return nil
}

func (m *memoryRepo) RawData() map[string]loans.LentBook {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return maps.Clone(m.lentBooks)
}

func (m *memoryRepo) FindLentBooks(ctx context.Context, at time.Time) ([]loans.LentBook, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	atUnix := uint64(at.Unix())
	result := make([]loans.LentBook, 0)
	for _, book := range m.lentBooks {
		if book.TakenAt <= atUnix && !(book.Returned && book.ReturnedAt <= atUnix) {
			result = append(result, book)
		}
	}
	return result, nil
}

func (m *memoryRepo) FindOverdueBooks(ctx context.Context, at time.Time) ([]loans.LentBook, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	atUnix := uint64(at.Unix())
	result := make([]loans.LentBook, 0)
	for _, book := range m.lentBooks {
		if book.ReturnDeadline <= atUnix && !(book.Returned && book.ReturnedAt <= atUnix) {
			result = append(result, book)
		}
	}
	return result, nil
}

func (m *memoryRepo) TakeBook(ctx context.Context, book *loans.LentBook, totalStock uint) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.lentBooks[book.ID]; ok {
		return fail.ErrCollision
	}

	remainingStock := int64(totalStock)

	for _, lentBook := range m.lentBooks {
		if lentBook.BookID == book.BookID {
			remainingStock -= 1
		}
	}

	if remainingStock < 0 {
		return fail.ErrNoStock
	}

	m.lentBooks[book.ID] = *book

	return nil
}

func (m *memoryRepo) ReturnBook(ctx context.Context, book *loans.LentBook) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldBook, ok := m.lentBooks[book.ID]
	if !ok {
		return fail.ErrNotFound
	}

	if oldBook.Returned {
		return fail.ErrCollision
	}

	if oldBook.UserID != book.UserID || oldBook.BookID != book.BookID {
		// TODO: Shouldn't actually happen ever. Maybe a different error?
		return fail.ErrCollision
	}

	m.lentBooks[book.ID] = *book

	return nil
}

func (m *memoryRepo) FindLoansOf(ctx context.Context, userID string, bookID string) ([]loans.LentBook, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]loans.LentBook, 0)
	for _, book := range m.lentBooks {
		if (userID == "" || book.UserID == userID) && (bookID == "" || book.BookID == bookID) {
			result = append(result, book)
		}
	}
	return result, nil
}
