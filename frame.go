package errors

import (
	"fmt"
	"io"
	"runtime"
)

type Frame struct {
	// Make room for three PCs: the one we were asked for, what it called,
	// and possibly a PC for skipPleaseUseCallersFrames. See:
	// https://go.googlesource.com/go/+/032678e0fb/src/runtime/extern.go#169
	f [3]uintptr
}

// Caller returns a Frame that describes a frame on the caller's stack.
// The argument skip is the number of frames to skip over.
// Caller(0) returns the frame for the caller of Caller.
func Caller(skip int) Frame {
	var f Frame
	runtime.Callers(skip+1, f.f[:])
	return f
}

func (f Frame) Info() FrameInfo {
	frames := runtime.CallersFrames(f.f[:])

	if _, ok := frames.Next(); !ok {
		return FrameInfo{}
	}
	fr, ok := frames.Next()
	if !ok {
		return FrameInfo{}
	}

	return FrameInfo{
		Function: fr.Function,
		File:     fr.File,
		Line:     fr.Line,
	}
}

type FrameInfo struct {
	Function string `json:"func,omitempty"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
}

func (f FrameInfo) FormatLine(w io.Writer) {
	if f.Function != "" {
		_, _ = io.WriteString(w, f.Function)
		_, _ = io.WriteString(w, "\n\t")
	}
	if f.File != "" {
		_, _ = io.WriteString(w, fmt.Sprintf("%s:%d\n", f.File, f.Line))
	}
}
