package errors

import (
	"fmt"
	"io"
)

type MetadataError interface {
	Metadata() map[string]interface{}
}

// WithCode decorate error with a code
//
// which can response to caller
func WithMetadata(err error, extra map[string]interface{}) error {
	if err == nil {
		return nil
	}

	return &withMetadata{
		cause: err,
		extra: extra,
	}
}

type withMetadata struct {
	cause error
	extra map[string]interface{}
}

func (w *withMetadata) Metadata() map[string]interface{} {
	return w.extra
}

func (w *withMetadata) Error() string {
	return w.cause.Error()
}

func (w *withMetadata) Cause() error { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withMetadata) Unwrap() error { return w.cause }

func (w *withMetadata) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			return
		}

		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.cause.Error())
	}
}

// Metadata get all etra information of the error
func Metadata(err error) map[string]interface{} {
	result := make(map[string]interface{})

	doUnwrap(err, func(err error) bool {
		if u, ok := err.(MetadataError); ok {
			extra := u.Metadata()

			for k, v := range extra {
				result[k] = v
			}
		}

		return false
	})

	return result
}
