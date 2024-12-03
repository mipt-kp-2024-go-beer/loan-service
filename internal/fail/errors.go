package fail

import "errors"

var (
	ErrNotFound  = new("object not found")
	ErrCollision = new("object already exists")
	ErrForbidden = new("insufficient permissions")
	ErrNoStock   = new("insufficient stock")
)

func new(desc string) error {
	return errors.New(desc)
}
