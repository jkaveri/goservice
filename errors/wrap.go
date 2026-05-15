package errors

import (
	stderrors "errors"
	"fmt"
)

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(StackError); !ok {
		err = &withStack{
			error: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		err: err,
		msg: message,
	}

	return err
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(StackError); !ok {
		err = &withStack{
			error: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		err: err,
		msg: fmt.Sprintf(format, args...),
	}

	return err
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// Unwrap returns nil if the Unwrap method returns []error.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}
