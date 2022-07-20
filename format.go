package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func FormatError(err error, f fmt.State, verb rune) {
	if err == nil {
		_, _ = io.WriteString(f, "<nil>")
	}

	for err != nil {
		if x, ok := err.(interface {
			formatSelf(w io.Writer, withFrame bool)
		}); ok {
			x.formatSelf(f, true)
		} else {
			_, _ = io.WriteString(f, err.Error())
			_, _ = io.WriteString(f, "\n")
		}

		err = Unwrap(err)
	}
}

func FormatJson(err error) ([]byte, error) {
	if err == nil {
		return json.Marshal(err)
	}

	var stack []jsonError

	for err != nil {
		var e jsonError

		if x, ok := err.(interface {
			formatSelf(w io.Writer, withFrame bool)
		}); ok {
			b := strings.Builder{}
			x.formatSelf(&b, false)
			e.Error = b.String()
		} else {
			e.Error = err.Error()
		}

		if x, ok := err.(interface {
			Frame() Frame
		}); ok {
			e.FrameInfo = x.Frame().Info()
		}

		stack = append(stack, e)

		err = Unwrap(err)
	}

	return json.Marshal(stack)
}

type jsonError struct {
	FrameInfo
	Error string `json:"error"`
}
