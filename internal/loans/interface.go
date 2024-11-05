package loans

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

type Service interface {
}

type Repo interface {
}
