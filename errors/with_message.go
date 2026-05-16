package errors

import (
	"errors"
	"fmt"
	"strings"
)

// MessageError marks an error that carries a short message meant for clients
// (HTTP/gRPC bodies, UI), distinct from the full text in Error() which still
// chains the underlying causes for operators and tests.
//
// For service logs and debugging, prefer fmt.Sprintf("%+v", err): the outer
// WithMessage prints its short message on the first line, then the wrapped
// chain with each layer's detail where implemented.
type MessageError interface {
	Message() string
	Error() string
}

// WithMessage wraps err with a user-facing message. The returned value's
// Error() includes both message and err; Message() returns only message.
// If err is nil, WithMessage returns nil. When HasStack(err) is false, a stack
// trace is captured at this call site (same policy as Wrap).
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}

	stacked := ensureStack(err)

	return &withMessage{
		err:      stacked,
		msg:      message,
		hasStack: stacked != nil,
	}
}

// WithMessagef is like WithMessage but formats message with fmt.Sprintf.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	stacked := ensureStack(err)

	return &withMessage{
		err:      stacked,
		msg:      fmt.Sprintf(format, args...),
		hasStack: stacked != nil,
	}
}

// Message returns the user-facing text to put in an API response (or similar).
// It walks err's Unwrap chain and joins each MessageError.Message() segment
// with ": ", outermost first. Layers that do not implement MessageError are
// skipped (for example the root stdlib or fundamental error), so callers only
// surface annotations added with WithMessage or WithMessagef. For nil err,
// Message returns an empty string.
func Message(err error) string {
	var (
		sb    strings.Builder
		first = true
	)

	for err := err; err != nil; err = errors.Unwrap(err) {
		if msgErr, ok := err.(MessageError); ok {
			if first {
				first = false
			} else {
				sb.WriteString(": ")
			}

			sb.WriteString(msgErr.Message())
		}
	}

	return sb.String()
}

type withMessage struct {
	err      error
	msg      string
	hasStack bool
}

func (w *withMessage) Message() string {
	return w.msg
}

// HasStack reports whether the wrapped chain had a stack when this message was
// added.
func (w *withMessage) HasStack() bool { return w.hasStack }

func (w *withMessage) Error() string { return fmt.Sprintf("%s: %s", w.msg, w.err.Error()) }

func (w *withMessage) Unwrap() error { return w.err }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s\n%+v", w.msg, w.Unwrap())
			return
		}

		fallthrough
	case 's':
		fmt.Fprintf(s, "%s: %s", w.msg, w.err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}
