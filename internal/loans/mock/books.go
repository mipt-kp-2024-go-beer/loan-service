package mock

import (
	"context"
	"fmt"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/books"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
)

func NewBooksConn() books.Connection {
	return &implBooksConn{}
}

type implBooksConn struct {
}

func (*implBooksConn) LookupBook(ctx context.Context, ID string) (*books.Book, error) {
	switch ID {
	case "bad-id":
		return nil, fmt.Errorf("%w: pretend missing book", fail.ErrBookService)
	case "single-book":
		return &books.Book{
			// Note: here and elsewhere, in mocks I use non-UUID IDs for simplicity
			ID:          "single-book",
			Title:       "The Bible",
			Author:      "God Almighty",
			Description: "lorem ipsum",
			TotalStock:  1,
		}, nil
	case "multi-book":
		return &books.Book{
			ID:          "multi-book",
			Title:       "The Bible",
			Author:      "God Almighty",
			Description: "lorem ipsum",
			TotalStock:  5,
		}, nil
	}
	panic("Unexpected request to mock book service!")
}
