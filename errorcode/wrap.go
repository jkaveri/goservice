package errorcode

import errors "github.com/jkaveri/goservice/errors"

// Wrap annotates err with the supplied message and tags it with the given
// error code. The returned error captures a stack trace at the call site and
// preserves err for use with [errors.Is] / [errors.As] / [errors.Unwrap].
// If err is nil, Wrap returns nil.
func Wrap(err error, code string, msg string) error {
	if err == nil {
		return nil
	}

	return errors.WithCode(errors.Wrap(err, msg), code)
}

// Wrapf is like [Wrap] but formats the annotation message with fmt.Sprintf
// using format and args. The returned error captures a stack trace at the
// call site and preserves err for use with [errors.Is] / [errors.As] /
// [errors.Unwrap]. If err is nil, Wrapf returns nil.
func Wrapf(err error, code string, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return errors.WithCode(errors.Wrapf(err, format, args...), code)
}
