package errors_test

import (
	"fmt"
	"github.com/jeanschmitt/errors"
)

func ExampleWrapWith() {
	err := errors.New("cause")

	// Wrapping the error with a custom, structured error type
	structured := errors.WrapWith(err, &HTTPError{code: 404})
	fmt.Println(structured)

	// Wrapping the error with a package-level error
	notFound := errors.WrapWith(structured, ErrNotFound)
	fmt.Println(notFound)

	// It is possible to determine if the cause is err
	fmt.Println(errors.Is(notFound, err))

	var httpError *HTTPError
	errors.As(notFound, &httpError)

	// It is possible to use errors.As to get one of the underlying errors struct
	fmt.Println(httpError.code)

	// Output:
	// http 404 error: cause
	// not found: http 404 error: cause
	// true
	// 404
}

var ErrNotFound = errors.New("not found")

type HTTPError struct {
	code int
}

func (e *HTTPError) Error() string { return fmt.Sprintf("http %d error", e.code) }
