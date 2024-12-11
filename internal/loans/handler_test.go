package loans_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/loans/mock"
)

func performRequest(t *testing.T, req *http.Request, useInternal bool) *httptest.ResponseRecorder {
	t.Helper()

	service := mock.NewService()

	router := chi.NewRouter()
	routerInternal := chi.NewRouter()

	h := loans.NewHandler(router, routerInternal, service)
	h.Register()

	rr := httptest.NewRecorder()

	if useInternal {
		routerInternal.ServeHTTP(rr, req)
	} else {
		router.ServeHTTP(rr, req)
	}

	return rr
}

// Public API

func TestPostBookTake(t *testing.T) {
	// POST /api/v1/book/{bookID}/take

	t.Run("basic", func(t *testing.T) {
		r, err := http.NewRequest(
			"POST",
			"/api/v1/book/good-book/take",
			strings.NewReader("auth=good-token"),
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, rr.Code)
		}
		if diff := cmp.Diff("{}\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad token", func(t *testing.T) {
		r, err := http.NewRequest(
			"POST",
			"/api/v1/book/good-book/take",
			strings.NewReader("auth=bad-token"),
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusForbidden {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusForbidden, rr.Code)
		}
		if diff := cmp.Diff("insufficient permissions\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad book", func(t *testing.T) {
		r, err := http.NewRequest(
			"POST",
			"/api/v1/book/bad-book/take",
			strings.NewReader("auth=good-token"),
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusNotFound {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusNotFound, rr.Code)
		}
		if diff := cmp.Diff("insufficient stock\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestPostBookReturn(t *testing.T) {
	// POST /api/v1/book/{bookID}/return

	t.Run("basic", func(t *testing.T) {
		r, err := http.NewRequest(
			"POST",
			"/api/v1/book/good-book/return",
			strings.NewReader("auth=good-token"),
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, rr.Code)
		}
		if diff := cmp.Diff("{}\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad token", func(t *testing.T) {
		r, err := http.NewRequest(
			"POST",
			"/api/v1/book/good-book/return",
			strings.NewReader("auth=bad-token"),
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusForbidden {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusForbidden, rr.Code)
		}
		if diff := cmp.Diff("insufficient permissions\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad book", func(t *testing.T) {
		r, err := http.NewRequest(
			"POST",
			"/api/v1/book/bad-book/return",
			strings.NewReader("auth=good-token"),
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusNotFound {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusNotFound, rr.Code)
		}
		if diff := cmp.Diff("object not found\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestGetBookAvailable(t *testing.T) {
	// GET /api/v1/book/{bookID}/avail

	t.Run("basic", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/book/good-book/avail?auth=good-token",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, rr.Code)
		}
		if diff := cmp.Diff("{\"available\":10}\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad token", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/book/good-book/avail?auth=bad-token",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusForbidden {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusForbidden, rr.Code)
		}
		if diff := cmp.Diff("insufficient permissions\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad book", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/book/bad-book/avail?auth=good-token",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusNotFound {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusNotFound, rr.Code)
		}
		if diff := cmp.Diff("object not found\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestGetReserved(t *testing.T) {
	// GET /api/v1/reserved

	t.Run("basic", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/reserved?auth=good-token",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, rr.Code)
		}
		if diff := cmp.Diff("{\"reserved\":[{\"id\":\"loan-id\",\"user_id\":\"user-id\",\"book_id\":\"book-id\",\"taken_at\":123,\"return_deadline\":456,\"returned\":false,\"returned_at\":0}]}\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("time", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/reserved?auth=good-token&atTime=12345",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, rr.Code)
		}
		if diff := cmp.Diff("{\"reserved\":[{\"id\":\"loan-id\",\"user_id\":\"user-id\",\"book_id\":\"book-id\",\"taken_at\":123,\"return_deadline\":456,\"returned\":false,\"returned_at\":0}]}\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad token", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/reserved?auth=bad-token",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusForbidden {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusForbidden, rr.Code)
		}
		if diff := cmp.Diff("insufficient permissions\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad time", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/reserved?auth=good-token&atTime=xxx",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusBadRequest, rr.Code)
		}
		if diff := cmp.Diff("missing required parameters: failed to parse atTime: strconv.ParseInt: parsing \"xxx\": invalid syntax\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestGetOverdue(t *testing.T) {
	// GET /api/v1/overdue

	t.Run("basic", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/overdue?auth=good-token",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, rr.Code)
		}
		if diff := cmp.Diff("{\"overdue\":[{\"id\":\"loan-id\",\"user_id\":\"user-id\",\"book_id\":\"book-id\",\"taken_at\":123,\"return_deadline\":456,\"returned\":false,\"returned_at\":0}]}\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("time", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/overdue?auth=good-token&atTime=12345",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, rr.Code)
		}
		if diff := cmp.Diff("{\"overdue\":[{\"id\":\"loan-id\",\"user_id\":\"user-id\",\"book_id\":\"book-id\",\"taken_at\":123,\"return_deadline\":456,\"returned\":false,\"returned_at\":0}]}\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad token", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/overdue?auth=bad-token",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusForbidden {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusForbidden, rr.Code)
		}
		if diff := cmp.Diff("insufficient permissions\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("bad time", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/overdue?auth=good-token&atTime=xxx",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, false)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusBadRequest, rr.Code)
		}
		if diff := cmp.Diff("missing required parameters: failed to parse atTime: strconv.ParseInt: parsing \"xxx\": invalid syntax\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})
}

// Internal API

func TestGetUserLoans(t *testing.T) {
	// GET /api/v1/userloans/{userID}

	t.Run("basic", func(t *testing.T) {
		r, err := http.NewRequest(
			"GET",
			"/api/v1/userloans/good-user",
			nil,
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := performRequest(t, r, true)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, rr.Code)
		}
		if diff := cmp.Diff("{\"unreturned\":123}\n", rr.Body.String()); diff != "" {
			t.Errorf("response body mismatch (-want +got):\n%s", diff)
		}
	})
}
