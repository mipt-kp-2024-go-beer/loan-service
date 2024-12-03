package fail

import "errors"

var (
	ErrNotFound  = new("object not found")
	ErrCollision = new("object already exists")
)

func new(desc string) error {
	return errors.New(desc)
}
