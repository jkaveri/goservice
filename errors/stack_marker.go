package errors

// stackMarker is implemented by errors package wrappers that record whether the
// unwrap chain already had a stack when the wrapper was constructed.
type stackMarker interface {
	HasStack() bool
}

// HasStack reports whether err already carries a stack trace in its chain.
// For errors built with Wrap, WithMessage, WithCode, or WithMetadata, the
// answer is cached on the outermost such wrapper (O(1)). Foreign errors are
// checked only
// on the value passed in (StackError or not).
func HasStack(err error) bool {
	return errHasStack(err)
}

// ensureStack attaches a stack at the current call site when err has none yet.
// If err is nil or already has a stack, err is returned unchanged.
func ensureStack(err error) error {
	if err == nil || errHasStack(err) {
		return err
	}

	return &withStack{
		error: err,
		stack: callers(),
	}
}

func errHasStack(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(StackError); ok {
		return true
	}

	if m, ok := err.(stackMarker); ok {
		return m.HasStack()
	}

	return false
}
