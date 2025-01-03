package loans

import (
	"context"
	"time"
)

// LentBook stores the information about a book being lent to a user
type LentBook struct {
	// ID is the UUID of the operation
	ID string `json:"id"`
	// UserID is the UUID of the user having taken the book
	UserID string `json:"user_id"`
	// BookID is the UUID of the book being lent
	BookID string `json:"book_id"`
	// TakenAt is the timestamp (UTC) when the book was taken
	TakenAt uint64 `json:"taken_at"`
	// ReturnDeadline is the timestamp (UTC) when the book should be returned
	ReturnDeadline uint64 `json:"return_deadline"`
	// Returned is true if the book is already returned
	Returned bool `json:"returned"`
	// ReturnedAt is the timestamp (UTC) when the book was returned, if it was already
	ReturnedAt uint64 `json:"returned_at"`
}

// Service is the interface for the business logic module of this microservice
type Service interface {
	// TakeBook records than a book is taken at the current date and time,
	// if it is in stock and the user has permission to take it.
	// If userID is not empty, the book is taken on behalf of the user with the given ID
	TakeBook(ctx context.Context, authToken string, userID string, bookID string) error

	// ReturnBook records than a book is returned at the current date and time,
	// if the user has permission to do so (is the one who had taken the book or a librarian).
	// If userID is not empty, the book is returned on behalf of the user with the given ID
	ReturnBook(ctx context.Context, authToken string, userID string, bookID string) error

	// CountAvailableBook returns the number of copies available for the given book,
	// if the user has permission to inquire this.
	CountAvailableBook(ctx context.Context, authToken string, bookID string) (uint, error)

	// ListReservations returns the list of books lent out at
	// the given time (now by default), if the user has permission to do so.
	ListReservations(ctx context.Context, authToken string, at time.Time) ([]LentBook, error)

	// ListOverdue returns the list of books that will become overdue by
	// the given time (now by default), if the user has permission to do so.
	ListOverdue(ctx context.Context, authToken string, at time.Time) ([]LentBook, error)

	// GetUserLoans returns how many unreturned lent books a particular user has at the moment
	GetUserLoans(ctx context.Context, userID string) (uint, error)

	// TODO: Some statistics? Clean up database?
}

// Repo is the interface for the memory module of this microservice
type Repo interface {
	// FindLentBooks returns the list of books lent out at the given time
	FindLentBooks(ctx context.Context, at time.Time) ([]LentBook, error)

	// FindOverdueBooks returns the list of books that will become overdue by the given time
	FindOverdueBooks(ctx context.Context, at time.Time) ([]LentBook, error)

	// TakeBook tests that the book isn't out of stock and registers it as taken.
	// book's fields must be set as if it was already taken
	TakeBook(ctx context.Context, book *LentBook, totalStock uint) error

	// ReturnBook tests that the book is taken and registers it as returned.
	// book's fields must be set as if it was already returned
	ReturnBook(ctx context.Context, book *LentBook) error

	// FindLoansOf finds all loans of a particular book by a particular user.
	// If either of (userID, bookID) is empty, that criterion is ignored
	FindLoansOf(ctx context.Context, userID string, bookID string) ([]LentBook, error)
}
