package errors

import (
	"fmt"
	"io"
)

// New returns an error containing the given message.
func New(message string) error {
	// The useful caller is New's caller
	return newFundamental(message, 1)
}

// Errorf returns an error containing the given formatted message.
func Errorf(format string, a ...any) error {
	// The useful caller is Errorf's caller
	return newFundamental(fmt.Sprintf(format, a...), 1)
}

// fundamental is an error only composed by a message.
type fundamental struct {
	msg   string
	frame Frame
}

func newFundamental(message string, skip int) error {
	return &fundamental{
		msg:   message,
		frame: Caller(skip + 1),
	}
}

func (e *fundamental) Error() string {
	return e.msg
}

func (e *fundamental) Format(f fmt.State, verb rune) { FormatError(e, f, verb) }

func (e *fundamental) Frame() Frame { return e.frame }

func (e *fundamental) formatSelf(w io.Writer, withFrame bool) {
	_, _ = io.WriteString(w, e.msg)
	if withFrame {
		_, _ = io.WriteString(w, ":\n")
		e.frame.Info().FormatLine(w)
	}
}

type Iser interface {
	Is(target error) bool
}

type Aser interface {
	As(target any) bool
}

type Unwrapper interface {
	Unwrap() error
}
