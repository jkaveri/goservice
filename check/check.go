package check

import (
	errors "github.com/jkaveri/goservice/errors"
)

// PanicIfError panic if error not nil.
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// PanicIfNil panic if val is nil.
func PanicIfNil(val interface{}, msg string, args ...interface{}) {
	PanicIf(val == nil, msg, args...)
}

// nolint
func PanicIf(condition bool, msg string, args ...interface{}) {
	if condition {
		panic(errors.Errorf(msg, args...))
	}
}

// nolint
func Ternary[T any](condition bool, trueResult T, falseResult T) T {
	if condition {
		return trueResult
	}

	return falseResult
}
