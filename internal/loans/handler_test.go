package loans_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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

	// TODO
}

func TestPostBookReturn(t *testing.T) {
	// POST /api/v1/book/{bookID}/return

	// TODO
}

func TestGetBookAvailable(t *testing.T) {
	// GET /api/v1/book/{bookID}/avail

	// TODO
}

func TestGetReserved(t *testing.T) {
	// GET /api/v1/reserved

	// TODO
}

func TestGetOverdue(t *testing.T) {
	// GET /api/v1/overdue

	// TODO
}

// Internal API

func TestGetUserLoans(t *testing.T) {
	// GET /api/v1/userloans/{userID}

	// TODO
}
