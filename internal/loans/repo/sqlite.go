package repo

import (
	"context"
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

	return result, rows.Err()
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

func (s *sqliteRepo) FindLentBooks(ctx context.Context, at time.Time) ([]loans.LentBook, error) {
	// TODO: Here and elsewhere, should I somehow respect context during mutex acquisition?
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rows, err := s.db.QueryContext(
		ctx,
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

func (s *sqliteRepo) FindOverdueBooks(ctx context.Context, at time.Time) ([]loans.LentBook, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rows, err := s.db.QueryContext(
		ctx,
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

func (s *sqliteRepo) TakeBook(ctx context.Context, book *loans.LentBook, totalStock uint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var lentStock uint
	err = tx.QueryRowContext(
		ctx,
		"SELECT count(*) FROM books WHERE id = ?",
		book.BookID,
	).Scan(&lentStock)
	if err != nil {
		return err
	}

	if lentStock >= totalStock {
		return fail.ErrNoStock
	}

	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO lent_books (id, user_id, book_id, taken_at, return_deadline, returned, returned_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		book.ID, book.UserID, book.BookID, book.TakenAt, book.ReturnDeadline, false, 0,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fail.ErrCollision
	}

	tx.Commit()

	return nil
}

func (s *sqliteRepo) ReturnBook(ctx context.Context, book *loans.LentBook) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	result, err := s.db.ExecContext(
		ctx,
		"UPDATE lent_books SET returned = TRUE, returned_at = ? WHERE id = ? AND returned = FALSE",
		book.ReturnedAt, book.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fail.ErrCollision
	}

	return nil
}

func (s *sqliteRepo) FindLoansOf(ctx context.Context, userID string, bookID string) ([]loans.LentBook, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rows, err := s.db.QueryContext(
		ctx,
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
