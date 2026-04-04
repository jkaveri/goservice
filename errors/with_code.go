package errors

import (
	"fmt"
	"io"
)

type CodeError interface {
	Code() string
}

// WithCode decorate error with a code
//
// which can response to caller
func WithCode(err error, code string) error {
	if err == nil {
		return nil
	}

	return &withCode{
		cause: err,
		code:  code,
	}
}

// Code get latest error code in the error chain
func Code(err error) string {
	code := ""

	doUnwrap(err, func(err error) bool {
		if u, ok := err.(CodeError); ok {
			code = u.Code()
			return true
		}

		return false
	})

	return code
}

// ContainsCode check to see error code in the chain or not?
func ContainsCode(err error, code string) bool {
	contains := false

	doUnwrap(err, func(err error) bool {
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
	cause error
	code  string
}

func (w *withCode) Code() string {
	return w.code
}

func (w *withCode) Error() string {
	return w.cause.Error()
}

func (w *withCode) Cause() error { return w.cause }

func (w *withCode) Unwrap() error {
	return w.cause
}

func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "[%s] %+v\n", w.code, w.Cause())
			return
		}

		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
