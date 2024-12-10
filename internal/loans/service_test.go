package loans_test

import (
	"context"
	"errors"
	"maps"
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans/mock"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans/repo"
)

const bookReturnDeadline = 48 * time.Hour

func makeService(t *testing.T) (context.Context, loans.Service, repo.TestMemoryRepo) {
	t.Helper()

	ctx := context.Background()

	userSvc := mock.NewUsersConn()
	bookSvc := mock.NewBooksConn()

	repo := repo.NewMemoryRepo("memory://")

	service := loans.NewService(repo, userSvc, bookSvc, bookReturnDeadline)

	return ctx, service, repo
}

func TestService_TakeBook(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		repo.ResetRawData(map[string]loans.LentBook{})

		err := service.TakeBook(ctx, "token-regular-user", "vasya-pupkin", "multi-book")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lentBooks := slices.Collect(maps.Values(repo.RawData()))
		if len(lentBooks) != 1 {
			t.Fatalf("expected 1 lent book, got %d", len(lentBooks))
		}

		got := lentBooks[0]
		want := loans.LentBook{
			ID:             got.ID,
			BookID:         "multi-book",
			UserID:         "vasya-pupkin",
			TakenAt:        got.TakenAt,
			ReturnDeadline: got.TakenAt + uint64(bookReturnDeadline.Seconds()),
			Returned:       false,
			ReturnedAt:     got.ReturnedAt,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("implicit user", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		repo.ResetRawData(map[string]loans.LentBook{})

		err := service.TakeBook(ctx, "token-regular-user", "", "multi-book")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lentBooks := slices.Collect(maps.Values(repo.RawData()))
		if len(lentBooks) != 1 {
			t.Fatalf("expected 1 lent book, got %d", len(lentBooks))
		}

		got := lentBooks[0]
		want := loans.LentBook{
			ID:             got.ID,
			BookID:         "multi-book",
			UserID:         "vasya-pupkin",
			TakenAt:        got.TakenAt,
			ReturnDeadline: got.TakenAt + uint64(bookReturnDeadline.Seconds()),
			Returned:       false,
			ReturnedAt:     got.ReturnedAt,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("other user", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		repo.ResetRawData(map[string]loans.LentBook{})

		err := service.TakeBook(ctx, "token-librarian", "vasya-pupkin", "multi-book")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lentBooks := slices.Collect(maps.Values(repo.RawData()))
		if len(lentBooks) != 1 {
			t.Fatalf("expected 1 lent book, got %d", len(lentBooks))
		}

		got := lentBooks[0]
		want := loans.LentBook{
			ID:             got.ID,
			BookID:         "multi-book",
			UserID:         "vasya-pupkin",
			TakenAt:        got.TakenAt,
			ReturnDeadline: got.TakenAt + uint64(bookReturnDeadline.Seconds()),
			Returned:       false,
			ReturnedAt:     got.ReturnedAt,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad permissions", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		repo.ResetRawData(map[string]loans.LentBook{})

		err := service.TakeBook(ctx, "token-regular-user", "yuuko-shirakawa", "multi-book")
		if !errors.Is(err, fail.ErrForbidden) {
			t.Errorf("wrong error: want %v, got %v", fail.ErrForbidden, err)
		}
	})

	t.Run("bad book", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		repo.ResetRawData(map[string]loans.LentBook{})

		err := service.TakeBook(ctx, "token-regular-user", "", "bad-id")
		if !errors.Is(err, fail.ErrBookService) {
			t.Errorf("wrong error: want %v, got %v", fail.ErrBookService, err)
		}
	})

	t.Run("no stock", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		repo.ResetRawData(map[string]loans.LentBook{
			"blah-blah-blah": {
				ID:             "blah-blah-blah",
				BookID:         "single-book",
				UserID:         "yuuko-shirakawa",
				TakenAt:        123,
				ReturnDeadline: 0xffffffff00000000,
				Returned:       false,
				ReturnedAt:     0,
			},
		})

		err := service.TakeBook(ctx, "token-regular-user", "", "single-book")
		if !errors.Is(err, fail.ErrNoStock) {
			t.Errorf("wrong error: want %v, got %v", fail.ErrNoStock, err)
		}
	})
}

func TestService_ReturnBook(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		bookPre := loans.LentBook{
			ID:             "blah-blah-blah",
			BookID:         "single-book",
			UserID:         "vasya-pupkin",
			TakenAt:        uint64(time.Now().Unix()) - 123,
			ReturnDeadline: uint64(time.Now().Unix()) + 123,
			Returned:       false,
			ReturnedAt:     0,
		}
		repo.ResetRawData(map[string]loans.LentBook{
			"blah-blah-blah": bookPre,
		})

		err := service.ReturnBook(ctx, "token-regular-user", "vasya-pupkin", "single-book")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lentBooks := slices.Collect(maps.Values(repo.RawData()))
		if len(lentBooks) != 1 {
			t.Fatalf("expected 1 lent book, got %d", len(lentBooks))
		}

		got := lentBooks[0]
		want := bookPre
		want.Returned = true
		want.ReturnedAt = got.ReturnedAt

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}

		upperBound := uint64(time.Now().Unix())
		lowerBound := upperBound - 5
		if !(lowerBound <= got.ReturnedAt && got.ReturnedAt <= upperBound) {
			t.Errorf("wrong returnedAt: want %d <= %d <= %d", lowerBound, got.ReturnedAt, upperBound)
		}
	})

	t.Run("implicit user", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		bookPre := loans.LentBook{
			ID:             "blah-blah-blah",
			BookID:         "single-book",
			UserID:         "vasya-pupkin",
			TakenAt:        uint64(time.Now().Unix()) - 123,
			ReturnDeadline: uint64(time.Now().Unix()) + 123,
			Returned:       false,
			ReturnedAt:     0,
		}
		repo.ResetRawData(map[string]loans.LentBook{
			"blah-blah-blah": bookPre,
		})

		err := service.ReturnBook(ctx, "token-regular-user", "", "single-book")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lentBooks := slices.Collect(maps.Values(repo.RawData()))
		if len(lentBooks) != 1 {
			t.Fatalf("expected 1 lent book, got %d", len(lentBooks))
		}

		got := lentBooks[0]
		want := bookPre
		want.Returned = true
		want.ReturnedAt = got.ReturnedAt

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}

		upperBound := uint64(time.Now().Unix())
		lowerBound := upperBound - 5
		if !(lowerBound <= got.ReturnedAt && got.ReturnedAt <= upperBound) {
			t.Errorf("wrong returnedAt: want %d <= %d <= %d", lowerBound, got.ReturnedAt, upperBound)
		}
	})

	t.Run("oldest lent", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		now := uint64(time.Now().Unix())
		booksPre := map[string]loans.LentBook{
			"pi-pi-pi": { // Non-returned, but different user
				ID:             "pi-pi-pi",
				BookID:         "multi-book",
				UserID:         "yuuko-shirakawa",
				TakenAt:        now - 100,
				ReturnDeadline: now,
				Returned:       false,
				ReturnedAt:     0,
			},
			"pa-pa-pa": { // Returned
				ID:             "pa-pa-pa",
				BookID:         "multi-book",
				UserID:         "vasya-pupkin",
				TakenAt:        now - 100,
				ReturnDeadline: now + 50,
				Returned:       true,
				ReturnedAt:     now - 1100,
			},
			"blah-blah-blah": { // Oldest non-returned
				ID:             "blah-blah-blah",
				BookID:         "multi-book",
				UserID:         "vasya-pupkin",
				TakenAt:        now - 100,
				ReturnDeadline: now + 100,
				Returned:       false,
				ReturnedAt:     0,
			},
			"pu-pu-pu": { // Non-returned, but newer
				ID:             "pu-pu-pu",
				BookID:         "multi-book",
				UserID:         "vasya-pupkin",
				TakenAt:        now - 100,
				ReturnDeadline: now + 200,
				Returned:       false,
				ReturnedAt:     0,
			},
		}
		repo.ResetRawData(booksPre)

		err := service.ReturnBook(ctx, "token-regular-user", "", "multi-book")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lentBooks := repo.RawData()
		if len(lentBooks) != 4 {
			t.Fatalf("expected 4 lent books, got %d", len(lentBooks))
		}

		got := lentBooks
		want := booksPre

		target := want["blah-blah-blah"]
		target.Returned = true
		target.ReturnedAt = got["blah-blah-blah"].ReturnedAt
		want["blah-blah-blah"] = target

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("result mismatch (-want +got):\n%s", diff)
		}

		upperBound := uint64(time.Now().Unix())
		lowerBound := upperBound - 5
		if !(lowerBound <= target.ReturnedAt && target.ReturnedAt <= upperBound) {
			t.Errorf("wrong returnedAt: want %d <= %d <= %d", lowerBound, target.ReturnedAt, upperBound)
		}
	})

	t.Run("not lent", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		repo.ResetRawData(map[string]loans.LentBook{})

		err := service.ReturnBook(ctx, "token-regular-user", "", "multi-book")
		if !errors.Is(err, fail.ErrNotFound) {
			t.Fatalf("wrong error: want %v, got %v", fail.ErrNotFound, err)
		}
	})

	t.Run("bad permissions", func(t *testing.T) {
		ctx, service, repo := makeService(t)
		bookPre := loans.LentBook{
			ID:             "blah-blah-blah",
			BookID:         "single-book",
			UserID:         "yuuko-shirakawa",
			TakenAt:        uint64(time.Now().Unix()) - 123,
			ReturnDeadline: uint64(time.Now().Unix()) + 123,
			Returned:       false,
			ReturnedAt:     0,
		}
		repo.ResetRawData(map[string]loans.LentBook{
			"blah-blah-blah": bookPre,
		})

		err := service.ReturnBook(ctx, "token-regular-user", "yuuko-shirakawa", "single-book")
		if !errors.Is(err, fail.ErrForbidden) {
			t.Fatalf("wrong error: want %v, got %v", fail.ErrForbidden, err)
		}
	})
}

func TestService_CountAvailableBook(t *testing.T) {
	ctx, service, repo := makeService(t)

	// TODO: Test serivce.CountAvailableBook
	_ = ctx
	_ = service
	_ = repo
}

func TestService_ListReservations(t *testing.T) {
	ctx, service, repo := makeService(t)

	// TODO: Test serivce.ListReservations
	_ = ctx
	_ = service
	_ = repo
}

func TestService_ListOverdue(t *testing.T) {
	ctx, service, repo := makeService(t)

	// TODO: Test serivce.ListOverdue
	_ = ctx
	_ = service
	_ = repo
}

func TestService_GetUserLoans(t *testing.T) {
	ctx, service, repo := makeService(t)

	// TODO: Test serivce.GetUserLoans
	_ = ctx
	_ = service
	_ = repo
}
