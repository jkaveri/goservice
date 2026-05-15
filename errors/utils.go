package errors

// WalkErrorChain visits err and its Unwrap() ancestors depth-first.
// For Unwrap() []error, each element is walked in order. If fn returns true,
// traversal stops.
// A nil err does not invoke fn.
func WalkErrorChain(err error, fn func(error) bool) {
	walkErrorChain(err, fn)
}

func walkErrorChain(err error, fn func(error) bool) bool {
	if err == nil {
		return false
	}

	if fn(err) {
		return true
	}

	switch e := err.(type) {
	case interface{ Unwrap() []error }:
		for _, item := range e.Unwrap() {
			if walkErrorChain(item, fn) {
				return true
			}
		}
	case interface{ Unwrap() error }:
		return walkErrorChain(e.Unwrap(), fn)
	}

	return false
}
