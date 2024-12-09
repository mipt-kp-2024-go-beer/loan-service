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

func (s *sqliteRepo) FindLentBooks(at time.Time) ([]loans.LentBook, error) {
	return nil, nil
}

func (s *sqliteRepo) FindOverdueBooks(at time.Time) ([]loans.LentBook, error) {
	return nil, nil
}

func (s *sqliteRepo) TakeBook(book loans.LentBook, totalStock uint) error {
	return nil
}

func (s *sqliteRepo) ReturnBook(book loans.LentBook) error {
	return nil
}

func (s *sqliteRepo) FindLoansOf(userID string, bookID string) ([]loans.LentBook, error) {
	return nil, nil
}
