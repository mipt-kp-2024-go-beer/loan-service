package repo

import (
	"time"

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
	return loans.LentBook{}, nil
}

func (m *memoryRepo) UpdateBook(book loans.LentBook) error {
	return nil
}

func (m *memoryRepo) InsertBook(book loans.LentBook) error {
	return nil
}

func (m *memoryRepo) FindLentBooks(at time.Time) ([]loans.LentBook, error) {
	return nil, nil
}

func (m *memoryRepo) FindOverdueBooks(at time.Time) ([]loans.LentBook, error) {
	return nil, nil
}
