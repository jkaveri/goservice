// Error handling helpers in this file follow patterns popularized by
// github.com/pkg/errors (BSD-2-Clause). Copyright (c) 2015, Dave Cheney.

package errors

import (
	"fmt"
	"io"
)

// New returns a leaf error whose Error() is message. It does not wrap another
// error; Unwrap on the result is nil. Prefer this constructor when you want
// values typed by this package rather than errors.New from the standard
// library.
func New(message string) error {
	return &fundamental{
		msg: message,
	}
}

// Errorf returns a leaf error whose Error() is fmt.Sprintf(format, args...).
// Like New, it does not wrap another error and Unwrap is nil.
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg: fmt.Sprintf(format, args...),
	}
}

// Wrap returns an error whose Unwrap() is err so errors.Is and errors.As
// traverse err from the returned value. When err is non-nil, Error() formats
// as "message: " plus err.Error(); when err is nil, Error() is message only
// (same as a leaf from New).
//
// When err is non-nil and the chain has no stack yet (see HasStack), Wrap
// captures a stack trace at this call site before adding message.
//
// Wrap does not treat nil err specially: the result is still non-nil with
// Unwrap returning nil. Callers who want nil when err is nil must check err
// before calling Wrap.
func Wrap(err error, message string) error {
	stacked := ensureStack(err)

	return &fundamental{
		cause:    stacked,
		msg:      message,
		hasStack: stacked != nil,
	}
}

// Wrapf is like Wrap but sets message using fmt.Sprintf(format, args...).
// The same nil err behavior as Wrap applies.
func Wrapf(err error, format string, args ...interface{}) error {
	stacked := ensureStack(err)

	return &fundamental{
		cause:    stacked,
		msg:      fmt.Sprintf(format, args...),
		hasStack: stacked != nil,
	}
}

// Unwrap returns the result of err's Unwrap method when err implements
// interface{ Unwrap() error }. Otherwise it returns nil. It matches the
// behavior of errors.Unwrap from the standard library for a single link.
func Unwrap(err error) error {
	if e, ok := err.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
	}

	return nil
}

// fundamental is the concrete type behind New, Errorf, Wrap, and Wrapf: a
// display message and an optional wrapped cause.
type fundamental struct {
	msg      string
	cause    error
	hasStack bool
}

// Unwrap returns the wrapped error, or nil for a leaf created with New or
// Errorf.
func (f *fundamental) Unwrap() error { return f.cause }

// HasStack reports whether the wrapped chain had a stack when this wrap was
// created.
func (f *fundamental) HasStack() bool { return f.hasStack }

// Error returns the message from construction when there is no wrapped error.
// When a cause is set and non-nil, Error returns "message: " concatenated with
// the cause's Error() string.
func (f *fundamental) Error() string {
	if f.cause == nil {
		return f.msg
	}

	return fmt.Sprintf("%s: %s", f.msg, f.cause.Error())
}

// Format implements fmt.Formatter. Verbs %s and plain %v match Error(). Verb
// %+v
// prints the message, then each wrapped cause on its own line with one extra
// tab of indent per nesting level (each line of the cause's %+v output is
// prefixed by a tab). Verb %q quotes the message only, not the causal chain.
func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)

			if f.cause != nil {
				fmt.Fprintf(s, "\n\t%+v", f.cause)
			}

			return
		}

		fallthrough
	case 's':
		if f.cause == nil {
			io.WriteString(s, f.msg)
			return
		}

		fmt.Fprintf(s, "%s: %s", f.msg, f.cause.Error())
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}
