package loans

import (
	"time"

	"github.com/google/uuid"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/books"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/users"
)

func NewService(repo Repo, users users.Connection, books books.Connection, returnDeadline time.Duration) Service {
	return &implService{
		repo:           repo,
		users:          users,
		books:          books,
		returnDeadline: returnDeadline,
	}
}

type implService struct {
	repo           Repo
	users          users.Connection
	books          books.Connection
	returnDeadline time.Duration
}

func (s *implService) TakeBook(authToken string, userID string, bookID string) error {
	user, err := s.users.VerifyToken(authToken)
	if err != nil {
		return err
	}

	if userID == "" {
		userID = user.ID
	}

	allowed := user.HasPerm(users.PermLoanBooks)
	if !allowed {
		return fail.ErrForbidden
	}

	book, err := s.books.LookupBook(bookID)
	if err != nil {
		return err
	}

	lentBook := LentBook{
		ID:             uuid.NewString(),
		UserID:         userID,
		BookID:         bookID,
		TakenAt:        uint64(time.Now().Unix()),
		ReturnDeadline: uint64(time.Now().Add(s.returnDeadline).Unix()),
		Returned:       false,
		ReturnedAt:     0,
	}

	err = s.repo.TakeBook(lentBook, book.TotalStock)
	return err
}

func (s *implService) ReturnBook(authToken string, userID string, bookID string) error {
	user, err := s.users.VerifyToken(authToken)
	if err != nil {
		return err
	}

	if userID == "" {
		userID = user.ID
	}

	allowed := user.HasPerm(users.PermLoanBooks)
	if !allowed {
		return fail.ErrForbidden
	}

	lentBooks, err := s.repo.FindLoansOf(userID, bookID)
	if err != nil {
		return err
	}

	if len(lentBooks) == 0 {
		return fail.ErrNotFound
	}

	oldestLentBook := lentBooks[0]
	for _, book := range lentBooks[1:] {
		if !book.Returned && book.ReturnDeadline < oldestLentBook.ReturnDeadline {
			oldestLentBook = book
		}
	}

	oldestLentBook.Returned = true
	oldestLentBook.ReturnedAt = uint64(time.Now().Unix())

	// Multiple DB operations without a common lock, but if a race condition
	// occurs (unlikely here), it will be detected as an error.
	err = s.repo.ReturnBook(oldestLentBook)
	return err
}

func (s *implService) CountAvailableBook(authToken string, bookID string) (uint, error) {
	user, err := s.users.VerifyToken(authToken)
	if err != nil {
		return 0, err
	}

	allowed := user.HasPerm(users.PermQueryAvailableStock)
	if !allowed {
		return 0, fail.ErrForbidden
	}

	lentBooks, err := s.repo.FindLoansOfBook(bookID)
	if err != nil {
		return 0, err
	}

	book, err := s.books.LookupBook(bookID)
	if err != nil {
		return 0, err
	}

	availableBooks := book.TotalStock
	for _, book := range lentBooks {
		if !book.Returned {
			availableBooks -= 1
		}
	}

	return availableBooks, nil
}

func (s *implService) ListReservations(authToken string, at time.Time) ([]LentBook, error) {
	user, err := s.users.VerifyToken(authToken)
	if err != nil {
		return nil, err
	}

	allowed := user.HasPerm(users.PermQueryReservations)
	if !allowed {
		return nil, fail.ErrForbidden
	}

	reservations, err := s.repo.FindLentBooks(at)
	return reservations, err
}

func (s *implService) ListOverdue(authToken string, at time.Time) ([]LentBook, error) {
	user, err := s.users.VerifyToken(authToken)
	if err != nil {
		return nil, err
	}

	allowed := user.HasPerm(users.PermQueryReservations)
	if !allowed {
		return nil, fail.ErrForbidden
	}

	overdue, err := s.repo.FindOverdueBooks(at)
	return overdue, err
}
