package errors

// Opaque returns an error with the same error formatting as err
// but that does not match err and cannot be unwrapped.
//
// It is useful to avoid leaking other packages errors.
func Opaque(err error) error {
	return opaque{err}
}

type opaque struct {
	error
}
