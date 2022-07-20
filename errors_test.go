package errors_test

import (
	"github.com/jeanschmitt/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	want := "error message"
	err := errors.New(want)
	require.Equal(t, want, err.Error())
}

func TestErrorf(t *testing.T) {
	want := "formatted \"error\""
	err := errors.Errorf("formatted %q", "error")
	require.Equal(t, want, err.Error())
}
