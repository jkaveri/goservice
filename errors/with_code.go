package errors

import (
	stderrors "errors"
	"fmt"
)

type CodeError interface {
	Code() string
	Error() string
}

// WithCode wraps err with a stable string code for callers and transport
// mapping.
// If err is nil, WithCode returns nil. When err has no stack yet, a stack trace
// is captured at this call site (see HasStack).
func WithCode(err error, code string) error {
	if err == nil {
		return nil
	}

	stacked := ensureStack(err)

	return &withCode{
		err:      stacked,
		code:     code,
		hasStack: stacked != nil,
	}
}

// Code returns the first CodeError code in err's chain, or "" if none.
func Code(err error) string {
	if codeErr, ok := stderrors.AsType[CodeError](err); ok {
		return codeErr.Code()
	}

	return ""
}

// ContainsCode check to see error code in the chain or not?
func ContainsCode(err error, code string) bool {
	contains := false

	WalkErrorChain(err, func(err error) bool {
		if u, ok := err.(CodeError); ok {
			if u.Code() == code {
				contains = true
				return true
			}
		}

		return false
	})

	return contains
}

type withCode struct {
	err      error
	code     string
	hasStack bool
}

func (w *withCode) Code() string {
	return w.code
}

// HasStack reports whether the wrapped chain had a stack when the code was
// attached.
func (w *withCode) HasStack() bool { return w.hasStack }

func (w *withCode) Error() string {
	return fmt.Sprintf("[%s] %s", w.code, w.err.Error())
}

func (w *withCode) Unwrap() error {
	return w.err
}

func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "[%s]\n\t%+v", w.code, w.Unwrap())
			return
		}

		fallthrough
	case 's':
		fmt.Fprintf(s, "[%s] %s", w.code, w.err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}
