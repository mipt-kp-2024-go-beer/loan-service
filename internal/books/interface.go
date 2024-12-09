package books

import "context"

// Book stores the information about a book
type Book struct {
	ID          string
	Title       string
	Author      string
	Description string
	TotalStock  uint
}

// Connection is the interface for the private API of the book microservice
type Connection interface {
	LookupBook(ctx context.Context, bookID string) (*Book, error)
}
