package errors

import "fmt"

// Wrap err with the given message and registers the stacktrace frame.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &wrapper{
		inner: err,
		outer: &fundamental{msg: message},
	}
}

// Wrapf err with the given formatted message and registers the stacktrace frame.
func Wrapf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}
	return &wrapper{
		inner: err,
		outer: &fundamental{msg: fmt.Sprintf(format, a...)},
	}
}

// WrapWith wraps err with a reason error.
//
// It is useful when wrapping a cause error with an already instanced error, or
// with a custom error.
func WrapWith(err error, reason error) error {
	if err == nil {
		return nil
	}
	if reason == nil {
		panic("cannot wrap with a nil error")
	}
	return &wrapper{
		inner: err,
		outer: reason,
	}
}

type wrapper struct {
	inner error
	outer error
}

func (w *wrapper) Error() string {
	return w.outer.Error() + ": " + w.inner.Error()
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
