package errors

import (
	"fmt"
	"io"
)

// Wrap err with the given message and registers the stacktrace frame.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return newWrapper(err, &fundamental{msg: message}, 1)
}

// Wrapf err with the given formatted message and registers the stacktrace frame.
func Wrapf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}
	return newWrapper(err, &fundamental{msg: fmt.Sprintf(format, a...)}, 1)
}

// WrapWith wraps err with another error.
//
// It is useful when wrapping a cause error with an already instanced error, or
// with a custom error.
func WrapWith(err error, with error) error {
	if err == nil {
		return nil
	}
	if with == nil {
		panic("cannot wrap with a nil error")
	}
	return newWrapper(err, with, 1)
}

type wrapper struct {
	inner error
	outer error
	frame Frame
}

func newWrapper(inner, outer error, skip int) error {
	return &wrapper{
		inner: inner,
		outer: outer,
		frame: Caller(skip + 1),
	}
}

func (w *wrapper) Error() string {
	return w.outer.Error() + ": " + w.inner.Error()
}

func (w *wrapper) Format(f fmt.State, verb rune) { FormatError(w, f, verb) }

func (w *wrapper) Frame() Frame { return w.frame }

func (w *wrapper) formatSelf(writer io.Writer, withFrame bool) {
	_, _ = io.WriteString(writer, w.outer.Error())
	if withFrame {
		_, _ = io.WriteString(writer, ":\n")
		w.frame.Info().FormatLine(writer)
	}
}

func (w *wrapper) Unwrap() error {
	return w.inner
}

// Is makes a non-recursive errors.Is call, because the next iteration must target
// w.inner.
func (w *wrapper) Is(target error) bool {
	return isNonRecursive(w.outer, target)
}

// As makes a non-recursive errors.As call, because the next iteration must target
// w.inner.
func (w *wrapper) As(target any) bool {
	return asNonRecursive(w.outer, target)
}
