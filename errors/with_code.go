package errors

import (
	stderrors "errors"
	"fmt"
	"io"
)

type CodeError interface {
	Code() string
	Error() string
}

// WithCode decorate error with a code
//
// which can response to caller
func WithCode(err error, code string) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:  err,
		code: code,
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
	err  error
	code string
}

func (w *withCode) Code() string {
	return w.code
}

func (w *withCode) Error() string {
	return w.err.Error()
}

func (w *withCode) Unwrap() error {
	return w.err
}

func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "[%s] %+v\n", w.code, w.Unwrap())
			return
		}

		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
