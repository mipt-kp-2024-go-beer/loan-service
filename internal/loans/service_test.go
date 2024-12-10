package loans_test

import (
	"context"
	"testing"
	"time"

	booksMock "github.com/mipt-kp-2024-go-beer/loan-service/internal/books/mock"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans/repo"
	usersMock "github.com/mipt-kp-2024-go-beer/loan-service/internal/users/mock"
)

/*
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
}
*/

const bookReturnDeadline = 48 * time.Hour

func TestService_TakeBook(t *testing.T) {
	ctx := context.Background()

	userSvc := usersMock.NewConn()
	bookSvc := booksMock.NewConn()

	repo := repo.NewMemoryRepo("memory://")

	// TODO: Set up repo

	service := loans.NewService(repo, userSvc, bookSvc, bookReturnDeadline)

	// TODO: Test serivce.TakeBook
	_ = ctx
	_ = service
}

func TestService_ReturnBook(t *testing.T) {
	ctx := context.Background()

	userSvc := usersMock.NewConn()
	bookSvc := booksMock.NewConn()

	repo := repo.NewMemoryRepo("memory://")

	// TODO: Set up repo

	service := loans.NewService(repo, userSvc, bookSvc, bookReturnDeadline)

	// TODO: Test serivce.ReturnBook
	_ = ctx
	_ = service
}

func TestService_CountAvailableBook(t *testing.T) {
	ctx := context.Background()

	userSvc := usersMock.NewConn()
	bookSvc := booksMock.NewConn()

	repo := repo.NewMemoryRepo("memory://")

	// TODO: Set up repo

	service := loans.NewService(repo, userSvc, bookSvc, bookReturnDeadline)

	// TODO: Test serivce.CountAvailableBook
	_ = ctx
	_ = service
}

func TestService_ListReservations(t *testing.T) {
	ctx := context.Background()

	userSvc := usersMock.NewConn()
	bookSvc := booksMock.NewConn()

	repo := repo.NewMemoryRepo("memory://")

	// TODO: Set up repo

	service := loans.NewService(repo, userSvc, bookSvc, bookReturnDeadline)

	// TODO: Test serivce.ListReservations
	_ = ctx
	_ = service
}

func TestService_ListOverdue(t *testing.T) {
	ctx := context.Background()

	userSvc := usersMock.NewConn()
	bookSvc := booksMock.NewConn()

	repo := repo.NewMemoryRepo("memory://")

	// TODO: Set up repo

	service := loans.NewService(repo, userSvc, bookSvc, bookReturnDeadline)

	// TODO: Test serivce.ListOverdue
	_ = ctx
	_ = service
}

func TestService_GetUserLoans(t *testing.T) {
	ctx := context.Background()

	userSvc := usersMock.NewConn()
	bookSvc := booksMock.NewConn()

	repo := repo.NewMemoryRepo("memory://")

	// TODO: Set up repo

	service := loans.NewService(repo, userSvc, bookSvc, bookReturnDeadline)

	// TODO: Test serivce.GetUserLoans
	_ = ctx
	_ = service
}
