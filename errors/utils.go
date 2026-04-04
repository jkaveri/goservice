package errors

// doUnwrap call fn(err) to handle the error
// if fn(err) return true, then stop the iteration
// if fn(err) return false, then continue the iteration
func doUnwrap(err error, fn func(err error) bool) {
	if err == nil {
		return
	}

	if fn(err) {
		return
	}

	switch unwrapped := err.(type) {
	case interface{ Unwrap() error }:
		doUnwrap(unwrapped.Unwrap(), fn)
	case interface{ Cause() error }:
		doUnwrap(unwrapped.Cause(), fn)
	case interface{ Unwrap() []error }:
		for _, e := range unwrapped.Unwrap() {
			doUnwrap(e, fn)
		}
	}
}
