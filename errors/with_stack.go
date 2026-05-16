package errors

import (
	"fmt"
	"io"
)

type StackError interface {
	StackTrace() StackTrace
	Error() string
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Unwrap() error { return w.error }

// HasStack is always true for errors that carry a captured stack.
func (w *withStack) HasStack() bool { return true }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Unwrap())
			w.stack.Format(s, verb)

			return
		}

		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

func (w *withStack) StackTrace() StackTrace {
	return w.stack.StackTrace()
}
