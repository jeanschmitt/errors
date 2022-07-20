package errors_test

import (
	"fmt"
	"github.com/jeanschmitt/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWrap_Nil(t *testing.T) {
	got := errors.Wrap(nil, "any message")
	require.Nil(t, got)
}

func TestWrap(t *testing.T) {
	cause := errors.New("cause")
	wrapped := errors.Wrap(cause, "middle")
	wrappedTwice := errors.Wrap(wrapped, "outer")

	require.Same(t, cause, errors.Unwrap(wrapped))
	require.Same(t, wrapped, errors.Unwrap(wrappedTwice))
	require.Same(t, cause, errors.Unwrap(errors.Unwrap(wrappedTwice)))

	require.Equal(t, "outer: middle: cause", wrappedTwice.Error())
}

func TestWrapf_Nil(t *testing.T) {
	got := errors.Wrapf(nil, "any %s message", "formatted")
	require.Nil(t, got)
}

func TestWrapf(t *testing.T) {
	cause := errors.New("cause")
	wrapped := errors.Wrapf(cause, "middle, %s formatted", "but")
	wrappedTwice := errors.Wrap(wrapped, "outer")

	require.Same(t, cause, errors.Unwrap(wrapped))
	require.Same(t, wrapped, errors.Unwrap(wrappedTwice))
	require.Same(t, cause, errors.Unwrap(errors.Unwrap(wrappedTwice)))

	require.Equal(t, "outer: middle, but formatted: cause", wrappedTwice.Error())
}

func TestWrapWith_NilErr(t *testing.T) {
	got := errors.WrapWith(nil, errors.New("any error"))
	require.Nil(t, got)
}

func TestWrapWith_NilWith(t *testing.T) {
	require.Panics(t, func() {
		_ = errors.WrapWith(errors.New("not nil"), nil)
	}, "wrapping with a nil error should panic")
}

func TestWrapper_Error(t *testing.T) {
	inner := errors.New("inner")
	outer := errors.New("outer")

	err := errors.WrapWith(inner, outer)
	require.Equal(t, "outer: inner", err.Error())
}

func TestWrapper_Unwrap(t *testing.T) {
	inner := errors.New("inner")
	err := errors.Wrap(inner, "outer")

	require.Implements(t, (*errors.Unwrapper)(nil), err)
	require.Same(t, inner, err.(errors.Unwrapper).Unwrap())
}

func TestWrapper_Is(t *testing.T) {
	inner := errors.New("inner")
	outerCause := errors.New("outer cause")
	outer := errors.Wrap(outerCause, "outer")

	err := errors.WrapWith(inner, outer)

	require.Implements(t, (*errors.Iser)(nil), err)
	require.True(t, err.(errors.Iser).Is(outer))
	require.False(t, err.(errors.Iser).Is(inner), "should only compare wrapper.outer")
	require.False(t, err.(errors.Iser).Is(outerCause), "should only compare with recursion")
}

func TestWrapper_As(t *testing.T) {
	inner := &otherTypedErr{code: 400}
	outerCause := &otherTypedErr{code: 500}
	outer := errors.WrapWith(outerCause, &typedErr{code: 600})

	err := errors.WrapWith(inner, outer)
	var target *typedErr
	var wrongTarget wrongTypeError
	var otherTarget *otherTypedErr

	require.Implements(t, (*errors.Aser)(nil), err)

	require.True(t, err.(errors.Aser).As(&target), "should success with a target of same type")
	require.Equal(t, 600, target.code)

	require.False(t, err.(errors.Aser).As(&wrongTarget), "should fail with a wrong target type")

	require.False(t, err.(errors.Aser).As(&otherTarget), "should only compare wrapper.outer, without recursion")
}

type typedErr struct {
	code int
}

func (e *typedErr) Error() string {
	return fmt.Sprintf("error %d", e.code)
}

type otherTypedErr struct {
	code int
}

func (e *otherTypedErr) Error() string {
	return fmt.Sprintf("error %d", e.code)
}

type wrongTypeError struct {
	error
}
