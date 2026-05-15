package errors

import (
	"fmt"
	"io"
)

// UserMessageError is implemented by errors that carry a client-facing message
// distinct from Error() (which may stay technical for logs).
//
// Deprecated: prefer API-layer client messages over error-wrapped copy.
type UserMessageError interface {
	// UserMessage returns text intended for API clients rather than logs.
	//
	// Deprecated: prefer API-layer client messages over error-wrapped copy.
	UserMessage() string
}

// WithUserMessage wraps err with a message safe to show to API clients.
// Error() still returns the underlying error string for logging and debugging.
// If userMessage is empty, err is returned unchanged.
//
// Deprecated: prefer API-layer client messages over error-wrapped copy.
func WithUserMessage(err error, userMessage string) error {
	if err == nil {
		return nil
	}

	if userMessage == "" {
		return err
	}

	return &withUserMessage{
		err:     err,
		message: userMessage,
	}
}

type withUserMessage struct {
	err     error
	message string
}

func (w *withUserMessage) UserMessage() string {
	return w.message
}

func (w *withUserMessage) Error() string {
	return w.err.Error()
}

func (w *withUserMessage) Unwrap() error { return w.err }

func (w *withUserMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(
				s,
				"user_message=%q %+v\n",
				w.message,
				w.Unwrap(),
			)

			return
		}

		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}

// GetUserMessage returns the first non-empty client message in err's chain, or
// "".
//
// Deprecated: prefer API-layer client messages over error-wrapped copy.
func GetUserMessage(err error) string {
	var msg string

	WalkErrorChain(err, func(e error) bool {
		if u, ok := e.(UserMessageError); ok {
			if m := u.UserMessage(); m != "" {
				msg = m
				return true
			}
		}

		return false
	})

	return msg
}
