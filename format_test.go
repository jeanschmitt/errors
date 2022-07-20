package errors_test

import (
	"fmt"
	"github.com/jeanschmitt/errors"
	"testing"
)

func TestFormat(t *testing.T) {
	err := wrapsPkg()

	fmt.Printf("%v\n", err)

	j, _ := errors.FormatJson(err)
	fmt.Println(string(j))
}

func cause() error {
	return errors.New("cause")
}

func wraps() error {
	return errors.Wrap(cause(), "wrapper message")
}

func wrapsPkg() error {
	return errors.WrapWith(wraps(), ErrPackage)
}

var ErrPackage = errors.New("package level error")
