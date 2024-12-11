package mock

import (
	"context"
	"time"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
)

func NewService() loans.Service {
	return &implService{}
}

type implService struct{}

func (s *implService) TakeBook(ctx context.Context, authToken string, userID string, bookID string) error {
	if authToken == "bad-token" {
		return fail.ErrForbidden
	}

	if bookID == "bad-book" {
		return fail.ErrNoStock
	}

	return nil
}

func (s *implService) ReturnBook(ctx context.Context, authToken string, userID string, bookID string) error {
	if authToken == "bad-token" {
		return fail.ErrForbidden
	}

	if bookID == "bad-book" {
		return fail.ErrNotFound
	}

	return nil
}

func (s *implService) CountAvailableBook(ctx context.Context, authToken string, bookID string) (uint, error) {
	if authToken == "bad-token" {
		return 0, fail.ErrForbidden
	}

	if bookID == "bad-book" {
		return 0, fail.ErrNotFound
	}

	return 10, nil
}

func (s *implService) ListReservations(ctx context.Context, authToken string, at time.Time) ([]loans.LentBook, error) {
	if authToken == "bad-token" {
		return nil, fail.ErrForbidden
	}

	return []loans.LentBook{
		{
			ID:             "loan-id",
			UserID:         "user-id",
			BookID:         "book-id",
			TakenAt:        123,
			ReturnDeadline: 456,
			Returned:       false,
			ReturnedAt:     0,
		},
	}, nil
}

func (s *implService) ListOverdue(ctx context.Context, authToken string, at time.Time) ([]loans.LentBook, error) {
	if authToken == "bad-token" {
		return nil, fail.ErrForbidden
	}

	return []loans.LentBook{
		{
			ID:             "loan-id",
			UserID:         "user-id",
			BookID:         "book-id",
			TakenAt:        123,
			ReturnDeadline: 456,
			Returned:       false,
			ReturnedAt:     0,
		},
	}, nil
}

func (s *implService) GetUserLoans(ctx context.Context, userID string) (uint, error) {
	return 123, nil
}
