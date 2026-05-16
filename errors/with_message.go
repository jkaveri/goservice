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

// NewMessage returns a MessageError with no wrapped cause. Error() and
// Message() are message. Prefer WithMessage when adding client text to an
// existing error chain.
func NewMessage(message string) error {
	return WithMessage(nil, message)
}

// NewMessagef is like NewMessage but formats message with fmt.Sprintf.
func NewMessagef(format string, args ...interface{}) error {
	return WithMessagef(nil, format, args...)
}

// WithMessage wraps err with a user-facing message. The returned value's
// Error() includes both message and err when err is non-nil; Message() returns
// only message. When err is nil, the result is a non-nil MessageError with
// Unwrap() == nil and Error() equal to message (user-facing text only, no
// internal chain). When HasStack(err) is false, a stack trace is captured at
// this call site (same policy as Wrap); a nil err has no stack until wrapped.
func WithMessage(err error, message string) error {
	stacked := ensureStack(err)

	return &withMessage{
		err:      stacked,
		msg:      message,
		hasStack: stacked != nil,
	}
}

// WithMessagef is like WithMessage but formats message with fmt.Sprintf.
func WithMessagef(err error, format string, args ...interface{}) error {
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
// surface annotations added with NewMessage, NewMessagef, WithMessage, or
// WithMessagef. For nil err,
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

func (w *withMessage) Error() string {
	if w.err == nil {
		return w.msg
	}

	return fmt.Sprintf("%s: %s", w.msg, w.err.Error())
}

func (w *withMessage) Unwrap() error { return w.err }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s", w.msg)

			if w.err != nil {
				fmt.Fprintf(s, "\n%+v", w.Unwrap())
			}

			return
		}

		fallthrough
	case 's':
		if w.err == nil {
			fmt.Fprintf(s, "%s", w.msg)
			return
		}

		fmt.Fprintf(s, "%s: %s", w.msg, w.err.Error())
	case 'q':
		if w.err == nil {
			fmt.Fprintf(s, "%q", w.msg)
			return
		}

		fmt.Fprintf(s, "%q", w.Error())
	}
}
