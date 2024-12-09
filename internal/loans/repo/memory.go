package repo

import (
	"sync"
	"time"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
)

func NewMemoryRepo(dsn string) loans.Repo {
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
	LookupBook(ID string) (loans.LentBook, error)
	UpdateBook(book loans.LentBook) error
	InsertBook(book loans.LentBook) error
}

func (m *memoryRepo) LookupBook(ID string) (loans.LentBook, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	book, ok := m.lentBooks[ID]
	if !ok {
		return loans.LentBook{}, fail.ErrNotFound
	}
	return book, nil
}

func (m *memoryRepo) UpdateBook(book loans.LentBook) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.lentBooks[book.ID]; !ok {
		return fail.ErrNotFound
	}
	m.lentBooks[book.ID] = book
	return nil
}

func (m *memoryRepo) InsertBook(book loans.LentBook) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.lentBooks[book.ID]; ok {
		return fail.ErrCollision
	}
	m.lentBooks[book.ID] = book
	return nil
}

func (m *memoryRepo) FindLentBooks(at time.Time) ([]loans.LentBook, error) {
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

func (m *memoryRepo) FindOverdueBooks(at time.Time) ([]loans.LentBook, error) {
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

func (m *memoryRepo) TakeBook(book loans.LentBook, totalStock uint) error {
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

	m.lentBooks[book.ID] = book

	return nil
}

func (m *memoryRepo) ReturnBook(book loans.LentBook) error {
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

	m.lentBooks[book.ID] = book

	return nil
}

func (m *memoryRepo) FindLoansOf(userID string, bookID string) ([]loans.LentBook, error) {
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
