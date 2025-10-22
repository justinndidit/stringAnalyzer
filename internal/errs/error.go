package errs

import (
	"errors"
	"fmt"
)

type InvalidTypeError struct {
	Field    string
	Expected string
	Got      string
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type for field '%s': expected %s, got %s", e.Field, e.Expected, e.Got)
}

var ErrNotFound = errors.New("string not found")
