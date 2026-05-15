package errors

import (
	"fmt"
	"io"
)

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}

	return &withMessage{
		err: err,
		msg: message,
	}
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withMessage{
		err: err,
		msg: fmt.Sprintf(format, args...),
	}
}

type withMessage struct {
	err error
	msg string
}

func (w *withMessage) Error() string { return w.msg }

func (w *withMessage) Unwrap() error { return w.err }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Unwrap())

			_, _ = io.WriteString(s, w.msg)

			return
		}

		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
