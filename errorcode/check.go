package errorcode

import errors "github.com/jkaveri/goservice/errors"

// IsErrorCode checks if the given error has the specified error code.
// Returns true if the error's code matches the provided code, false otherwise.
// If err is nil or doesn't have a code, it returns false.
func IsErrorCode(err error, code string) bool {
	return code == errors.Code(err)
}

// IsNotFound checks if the given error is a "not found" error.
// Returns true if the error's code matches CodeNotFound, false otherwise.
func IsNotFound(err error) bool {
	return CodeNotFound == errors.Code(err)
}

// IsUnauthorized checks if the given error is an "unauthorized" error.
// Returns true if the error's code matches CodeUnauthorized, false otherwise.
func IsUnauthorized(err error) bool {
	return CodeUnauthorized == errors.Code(err)
}

// IsDuplicated checks if the given error is a "duplicated" error.
// Returns true if the error's code matches CodeDuplicated, false otherwise.
func IsDuplicated(err error) bool {
	return CodeDuplicated == errors.Code(err)
}

// IsInternalServer checks if the given error is an "internal server" error.
// Returns true if the error's code matches CodeInternalServer, false otherwise.
func IsInternalServer(err error) bool {
	return CodeInternalServer == errors.Code(err)
}

// IsInvalidRequest checks if the given error is an "invalid request" error.
// Returns true if the error's code matches CodeInvalidRequest, false otherwise.
func IsInvalidRequest(err error) bool {
	return CodeInvalidRequest == errors.Code(err)
}

// IsNotAuthenticated checks if the given error is a "not authenticated" error.
// Returns true if the error's code matches CodeNotAuthenticated, false
// otherwise.
func IsNotAuthenticated(err error) bool {
	return CodeNotAuthenticated == errors.Code(err)
}
