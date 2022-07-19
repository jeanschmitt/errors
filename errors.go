package errors

import "fmt"

func New(message string) error {
	return &fundamental{msg: message}
}

func Errorf(format string, a ...any) error {
	return New(fmt.Sprintf(format, a...))
}

// fundamental is an error only composed by a message.
type fundamental struct {
	msg string
}

func (e *fundamental) Error() string {
	return e.msg
}
