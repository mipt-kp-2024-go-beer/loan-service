package repo

import (
	"time"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
)

func NewSqliteRepo() loans.Repo {
	return &sqliteRepo{}
}

type sqliteRepo struct {
}

func (m *sqliteRepo) LookupBook(bookID string) (loans.LentBook, error) {
	return loans.LentBook{}, nil
}

func (m *sqliteRepo) UpdateBook(book loans.LentBook) error {
	return nil
}

func (m *sqliteRepo) InsertBook(book loans.LentBook) error {
	return nil
}

func (m *sqliteRepo) FindLentBooks(at time.Time) ([]loans.LentBook, error) {
	return nil, nil
}

func (m *sqliteRepo) FindOverdueBooks(at time.Time) ([]loans.LentBook, error) {
	return nil, nil
}
