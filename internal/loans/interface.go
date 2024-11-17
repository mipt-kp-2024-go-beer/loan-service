package loans

import "time"

// LentBook stores the information about a book being lent to a user
type LentBook struct {
	// ID is the UUID of the operation
	ID string
	// UserID is the UUID of the user having taken the book
	UserID string
	// BookID is the UUID of the book being lent
	BookIDs string
	// TakenAt is the timestamp (UTC) when the book was taken
	TakenAt uint64
	// ReturnDeadline is the timestamp (UTC) when the book should be returned
	ReturnDeadline uint64
	// Returned is true if the book is already returned
	Returned bool
	// ReturnedAt is the timestamp (UTC) when the book was returned, if it was already
	ReturnedAt uint64
}

// Service is the interface for the business logic module of this microservice
type Service interface {
	// TakeBook records than a book is taken at the current date and time,
	// if it is in stock and the user has permission to take it.
	TakeBook(authToken string, bookID string) error

	// ReturnBook records than a book is returned at the current date and time,
	// if the user has permission to do so (is the one who had taken the book or a librarian).
	ReturnBook(authToken string, bookID string) error

	// CountAvailableBook returns the number of copies available for the given book,
	// if the user has permission to inquire this.
	CountAvailableBook(authToken string, bookID string) (uint, error)

	// ListReservations returns the list of books lent out at
	// the given time (now by default), if the user has permission to do so.
	ListReservations(authToken string, at time.Time) ([]LentBook, error)

	// ListOverdue returns the list of books that will become overdue by
	// the given time (now by default), if the user has permission to do so.
	ListOverdue(authToken string, at time.Time) ([]LentBook, error)

	// TODO: Some statistics? Clean up database?
}

// Repo is the interface for the memory module of this microservice
type Repo interface {
	// LookupBook looks up a LentBook by ID
	LookupBook(bookID string) (LentBook, error)

	// UpdateBook updates a LentBook with the given ID in the database
	UpdateBook(book LentBook) error

	// InsertBook adds a new LentBook to the database
	InsertBook(book LentBook) error

	// FindLentBooks returns the list of books lent out at the given time
	FindLentBooks(at time.Time) ([]LentBook, error)

	// FindOverdueBooks returns the list of books that will become overdue by the given time
	FindOverdueBooks(at time.Time) ([]LentBook, error)
}
