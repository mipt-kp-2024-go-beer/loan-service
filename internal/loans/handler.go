package loans

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	router  *chi.Mux
	service Service
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
		r.Get("/api/v1/book/{bookID}/avail", h.getBookStatus)

		r.Get("/api/v1/reserved", h.getReserved)
		r.Get("/api/v1/overdue", h.getOverdue)
	})
}

func (h *Handler) postBookTake(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) postBookReturn(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) getBookStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) getReserved(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) getOverdue(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
