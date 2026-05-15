package errors

import (
	stderrors "errors"
	"fmt"
	"io"
)

type MetadataError interface {
	Metadata() map[string]any
	Error() string
}

// WithCode decorate error with a code
//
// which can response to caller
func WithMetadata(err error, extra map[string]any) error {
	if err == nil {
		return nil
	}

	return &withMetadata{
		err:   err,
		extra: extra,
	}
}

type withMetadata struct {
	err   error
	extra map[string]any
}

func (w *withMetadata) Metadata() map[string]any {
	return w.extra
}

func (w *withMetadata) Error() string {
	return w.err.Error()
}

func (w *withMetadata) Unwrap() error { return w.err }

func (w *withMetadata) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Unwrap())
			return
		}

		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.err.Error())
	}
}

// Metadata get all etra information of the error
func Metadata(err error) map[string]any {
	if metadataErr, ok := stderrors.AsType[MetadataError](err); ok {
		return metadataErr.Metadata()
	}

	return nil
}
