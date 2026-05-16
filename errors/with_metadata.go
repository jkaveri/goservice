package errors

import (
	stderrors "errors"
	"fmt"
	"sort"
)

type MetadataError interface {
	Metadata() map[string]any
	Error() string
}

// WithMetadata wraps err with a shallow copy of extra for inspection (see
// Metadata). Error() still returns err.Error(); fmt.Sprintf("%+v", err) prints
// the wrapped chain and then a metadata block with sorted keys for stable logs.
func WithMetadata(err error, extra map[string]any) error {
	if err == nil {
		return nil
	}

	return &withMetadata{
		err:      err,
		extra:    extra,
		hasStack: errHasStack(err),
	}
}

type withMetadata struct {
	err      error
	extra    map[string]any
	hasStack bool
}

func (w *withMetadata) Metadata() map[string]any {
	return w.extra
}

// HasStack reports whether the wrapped chain had a stack when metadata was
// attached.
func (w *withMetadata) HasStack() bool { return w.hasStack }

func (w *withMetadata) Error() string {
	return w.err.Error()
}

func (w *withMetadata) Unwrap() error { return w.err }

func (w *withMetadata) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Unwrap())

			if len(w.extra) > 0 {
				fmt.Fprintf(s, "\nmetadata:")

				keys := make([]string, 0, len(w.extra))
				for k := range w.extra {
					keys = append(keys, k)
				}

				sort.Strings(keys)

				for _, k := range keys {
					fmt.Fprintf(s, "\t%s: %v", k, w.extra[k])
				}
			}

			return
		}

		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", w.err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// Metadata returns attached metadata, or nil when err does not implement
// MetadataError.
func Metadata(err error) map[string]any {
	if metadataErr, ok := stderrors.AsType[MetadataError](err); ok {
		return metadataErr.Metadata()
	}

	return nil
}
