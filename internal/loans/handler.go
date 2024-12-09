package loans

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
)

type Handler struct {
	router         *chi.Mux
	routerInternal *chi.Mux
	service        Service
}

func NewHandler(service Service) Handler {
	return Handler{
		router:  chi.NewRouter(),
		service: service,
	}
}

func (h *Handler) Register() {
	h.router.Group(func(r chi.Router) {
		r.Post("/api/v1/book/{bookID}/take", h.postBookTake)
		r.Post("/api/v1/book/{bookID}/return", h.postBookReturn)
		r.Get("/api/v1/book/{bookID}/avail", h.getBookAvailable)

		r.Get("/api/v1/reserved", h.getReserved)
		r.Get("/api/v1/overdue", h.getOverdue)
	})

	h.routerInternal.Group(func(r chi.Router) {
		r.Get("/api/v1/userloans/{userID}", h.getUserLoans)
	})
}

func writeJSONSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct{}{})
}

// Public API

func (h *Handler) postBookTake(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fail.WriteJSONError(w, fmt.Errorf("%w: %w", fail.ErrMissingParams, err))
		return
	}
	authToken := r.Form.Get("auth")
	userID := r.Form.Get("user")
	bookID := chi.URLParam(r, "bookID")
	if authToken == "" || bookID == "" {
		fail.WriteJSONError(w, fmt.Errorf("%w: %q", fail.ErrMissingParams, "auth, bookID"))
		return
	}

	err = h.service.TakeBook(authToken, userID, bookID)
	if err != nil {
		fail.WriteJSONError(w, err)
		return
	}

	writeJSONSuccess(w)
}

func (h *Handler) postBookReturn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fail.WriteJSONError(w, fmt.Errorf("%w: %w", fail.ErrMissingParams, err))
		return
	}
	authToken := r.Form.Get("auth")
	userID := r.Form.Get("user")
	bookID := chi.URLParam(r, "bookID")
	if authToken == "" || bookID == "" {
		fail.WriteJSONError(w, fmt.Errorf("%w: %q", fail.ErrMissingParams, "auth, bookID"))
		return
	}

	err = h.service.ReturnBook(authToken, userID, bookID)
	if err != nil {
		fail.WriteJSONError(w, err)
		return
	}

	writeJSONSuccess(w)
}

func (h *Handler) getBookAvailable(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fail.WriteJSONError(w, fmt.Errorf("%w: %w", fail.ErrMissingParams, err))
		return
	}
	authToken := r.Form.Get("auth")
	bookID := chi.URLParam(r, "bookID")
	if authToken == "" || bookID == "" {
		fail.WriteJSONError(w, fmt.Errorf("%w: %q", fail.ErrMissingParams, "auth, bookID"))
		return
	}

	available, err := h.service.CountAvailableBook(authToken, bookID)
	if err != nil {
		fail.WriteJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Available uint `json:"available"`
	}{
		Available: available,
	})
}

func (h *Handler) getReserved(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fail.WriteJSONError(w, fmt.Errorf("%w: %w", fail.ErrMissingParams, err))
		return
	}
	authToken := r.Form.Get("auth")
	atTime, err := strconv.ParseInt(r.Form.Get("atTime"), 10, 64)
	if err != nil {
		fail.WriteJSONError(w, fmt.Errorf("%w: failed to parse atTime: %w", fail.ErrMissingParams, err))
		return
	}
	if authToken == "" {
		fail.WriteJSONError(w, fmt.Errorf("%w: %q", fail.ErrMissingParams, "auth"))
		return
	}

	reserved, err := h.service.ListReservations(authToken, time.Unix(atTime, 0))
	if err != nil {
		fail.WriteJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Reserved []LentBook `json:"reserved"`
	}{
		Reserved: reserved,
	})
}

func (h *Handler) getOverdue(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fail.WriteJSONError(w, fmt.Errorf("%w: %w", fail.ErrMissingParams, err))
		return
	}
	authToken := r.Form.Get("auth")
	atTime, err := strconv.ParseInt(r.Form.Get("atTime"), 10, 64)
	if err != nil {
		fail.WriteJSONError(w, fmt.Errorf("%w: failed to parse atTime: %w", fail.ErrMissingParams, err))
		return
	}
	if authToken == "" {
		fail.WriteJSONError(w, fmt.Errorf("%w: %q", fail.ErrMissingParams, "auth"))
		return
	}

	overdue, err := h.service.ListOverdue(authToken, time.Unix(atTime, 0))
	if err != nil {
		fail.WriteJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Overdue []LentBook `json:"overdue"`
	}{
		Overdue: overdue,
	})
}

// Internal API

func (h *Handler) getUserLoans(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		fail.WriteJSONError(w, fmt.Errorf("%w: %q", fail.ErrMissingParams, "userID"))
		return
	}

	unreturned, err := h.service.GetUserLoans(userID)
	if err != nil {
		fail.WriteJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Unreturned uint `json:"unreturned"`
	}{
		Unreturned: unreturned,
	})
}
