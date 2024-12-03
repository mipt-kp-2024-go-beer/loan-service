package repo

import (
	"time"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
)

func NewMemoryRepo() loans.Repo {
	return &memoryRepo{
		make(map[string]loans.LentBook),
	}
}

type memoryRepo struct {
	lentBooks map[string]loans.LentBook
}

func (m *memoryRepo) LookupBook(bookID string) (loans.LentBook, error) {
	book, ok := m.lentBooks[bookID]
	if !ok {
		return loans.LentBook{}, fail.ErrNotFound
	}
	return book, nil
}

func (m *memoryRepo) UpdateBook(book loans.LentBook) error {
	if _, ok := m.lentBooks[book.ID]; !ok {
		return fail.ErrNotFound
	}
	m.lentBooks[book.ID] = book
	return nil
}

func (m *memoryRepo) InsertBook(book loans.LentBook) error {
	if _, ok := m.lentBooks[book.ID]; ok {
		return fail.ErrCollision
	}
	m.lentBooks[book.ID] = book
	return nil
}

func (m *memoryRepo) FindLentBooks(at time.Time) ([]loans.LentBook, error) {
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
	atUnix := uint64(at.Unix())
	result := make([]loans.LentBook, 0)
	for _, book := range m.lentBooks {
		if book.ReturnDeadline <= atUnix && !(book.Returned && book.ReturnedAt <= atUnix) {
			result = append(result, book)
		}
	}
	return result, nil
}
