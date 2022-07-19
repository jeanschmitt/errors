package errors_test

import (
	"fmt"
	"github.com/jeanschmitt/errors"
	"testing"
)

func TestAsGeneric(t *testing.T) {
	err := errors.New("not found")
	s := errors.WrapWith(err, &StructuredError{code: 404})
	wrapped := errors.Wrap(s, "wrapper")

	structured, ok := errors.AsGeneric[*StructuredError](wrapped)
	if !ok {
		t.FailNow()
	}

	if structured.code != 404 {
		t.FailNow()
	}
}

type StructuredError struct {
	code int
}

func (e *StructuredError) Error() string { return fmt.Sprintf("error %d", e.code) }
