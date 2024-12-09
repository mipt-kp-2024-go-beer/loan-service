package repo

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
)

func NewSqliteRepo(dsn string) (loans.Repo, error) {
	dsn, ok := strings.CutPrefix(dsn, "sqlite:")
	if !ok {
		return nil, fail.ErrInvalidDSN
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	return &sqliteRepo{
		mutex: sync.RWMutex{},
		db:    db,
	}, nil
}

type sqliteRepo struct {
	mutex sync.RWMutex
	db    *sql.DB
}

type sqliteLentBook struct {
	ID             sql.NullString
	UserID         sql.NullString
	BookID         sql.NullString
	TakenAt        sql.NullInt64
	ReturnDeadline sql.NullInt64
	Returned       sql.NullBool
	ReturnedAt     sql.NullInt64
}

func convertSqliteToReal(sqliteLentBook sqliteLentBook) (loans.LentBook, error) {
	if !(sqliteLentBook.ID.Valid &&
		sqliteLentBook.UserID.Valid &&
		sqliteLentBook.BookID.Valid &&
		sqliteLentBook.TakenAt.Valid &&
		sqliteLentBook.ReturnDeadline.Valid &&
		sqliteLentBook.Returned.Valid &&
		sqliteLentBook.ReturnedAt.Valid) {
		return loans.LentBook{}, fail.ErrMalformedStorage
	}

	return loans.LentBook{
		ID:             sqliteLentBook.ID.String,
		UserID:         sqliteLentBook.UserID.String,
		BookID:         sqliteLentBook.BookID.String,
		TakenAt:        uint64(sqliteLentBook.TakenAt.Int64),
		ReturnDeadline: uint64(sqliteLentBook.ReturnDeadline.Int64),
		Returned:       sqliteLentBook.Returned.Bool,
		ReturnedAt:     uint64(sqliteLentBook.ReturnedAt.Int64),
	}, nil
}

func convertRowsToReal(rows *sql.Rows) ([]loans.LentBook, error) {
	result := make([]loans.LentBook, 0)

	for rows.Next() {
		var sqliteLentBook sqliteLentBook
		if err := rows.Scan(&sqliteLentBook); err != nil {
			return nil, err
		}

		realLentBook, err := convertSqliteToReal(sqliteLentBook)
		if err != nil {
			return nil, err
		}

		result = append(result, realLentBook)
	}

	return result, nil
}

func convertRealToSqlite(realLentBook loans.LentBook) sqliteLentBook {
	return sqliteLentBook{
		ID:             sql.NullString{String: realLentBook.ID, Valid: true},
		UserID:         sql.NullString{String: realLentBook.UserID, Valid: true},
		BookID:         sql.NullString{String: realLentBook.BookID, Valid: true},
		TakenAt:        sql.NullInt64{Int64: int64(realLentBook.TakenAt), Valid: true},
		ReturnDeadline: sql.NullInt64{Int64: int64(realLentBook.ReturnDeadline), Valid: true},
		Returned:       sql.NullBool{Bool: realLentBook.Returned, Valid: true},
		ReturnedAt:     sql.NullInt64{Int64: int64(realLentBook.ReturnedAt), Valid: true},
	}
}

func (s *sqliteRepo) FindLentBooks(at time.Time) ([]loans.LentBook, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rows, err := s.db.Query(
		"SELECT * FROM lent_books WHERE taken_at <= ? AND NOT (returned AND returned_at <= ?)",
		at.Unix(), at.Unix(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := convertRowsToReal(rows)
	return result, err
}

func (s *sqliteRepo) FindOverdueBooks(at time.Time) ([]loans.LentBook, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rows, err := s.db.Query(
		"SELECT * FROM lent_books WHERE return_deadline <= ? AND NOT (returned AND returned_at <= ?)",
		at.Unix(), at.Unix(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := convertRowsToReal(rows)
	return result, err
}

func (s *sqliteRepo) TakeBook(book loans.LentBook, totalStock uint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return nil
}

func (s *sqliteRepo) ReturnBook(book loans.LentBook) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return nil
}

func (s *sqliteRepo) FindLoansOf(userID string, bookID string) ([]loans.LentBook, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rows, err := s.db.Query(
		"SELECT * FROM lent_books WHERE (? OR user_id = ?) AND (? OR book_id = ?)",
		userID == "", userID, bookID == "", bookID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := convertRowsToReal(rows)
	return result, err
}
