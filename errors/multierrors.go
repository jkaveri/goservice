package errors

import (
	stderrors "errors"
)

// Join create and error which contains all none nil errors
//
// this function returns nil when all errors are nil
func Join(errs ...error) error {
	return stderrors.Join(errs...)
}
