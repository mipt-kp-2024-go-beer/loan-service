package fail

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound         = new("object not found")
	ErrCollision        = new("object already exists")
	ErrForbidden        = new("insufficient permissions")
	ErrNoStock          = new("insufficient stock")
	ErrMissingParams    = new("missing required parameters")
	ErrInvalidDSN       = new("unrecognized data source name")
	ErrMalformedStorage = new("malformed storage")
	ErrUserService      = new("user service error")
	ErrBookService      = new("book service error")
)

func new(desc string) error {
	return errors.New(desc)
}

// HTTPErrorCode returns the HTTP error code for the given error.
func HTTPErrorCode(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrCollision):
		return http.StatusConflict
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrMissingParams):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// WriteError writes the right error code and the error message to the given writer.
func WriteError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), HTTPErrorCode(err))
}
