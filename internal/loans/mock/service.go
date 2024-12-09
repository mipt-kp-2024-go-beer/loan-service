package mock

import (
	"context"
	"time"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
)

func NewService() loans.Service {
	return &implService{}
}

type implService struct{}

func (s *implService) TakeBook(ctx context.Context, authToken string, userID string, bookID string) error {
	// TODO

	return nil
}

func (s *implService) ReturnBook(ctx context.Context, authToken string, userID string, bookID string) error {
	// TODO

	return nil
}

func (s *implService) CountAvailableBook(ctx context.Context, authToken string, bookID string) (uint, error) {
	// TODO

	return 0, nil
}

func (s *implService) ListReservations(ctx context.Context, authToken string, at time.Time) ([]loans.LentBook, error) {
	// TODO

	return nil, nil
}

func (s *implService) ListOverdue(ctx context.Context, authToken string, at time.Time) ([]loans.LentBook, error) {
	// TODO

	return nil, nil
}

func (s *implService) GetUserLoans(ctx context.Context, userID string) (uint, error) {
	// TODO

	return 0, nil
}
