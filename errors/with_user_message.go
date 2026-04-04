package errors

import (
	"fmt"
	"io"
)

// UserMessageError is implemented by errors that carry a client-facing message
// distinct from Error() (which may stay technical for logs).
type UserMessageError interface {
	UserMessage() string
}

// WithUserMessage wraps err with a message safe to show to API clients.
// Error() still returns the underlying error string for logging and debugging.
// If userMessage is empty, err is returned unchanged.
func WithUserMessage(err error, userMessage string) error {
	if err == nil {
		return nil
	}

	if userMessage == "" {
		return err
	}

	return &withUserMessage{
		cause:   err,
		message: userMessage,
	}
}

type withUserMessage struct {
	cause   error
	message string
}

func (w *withUserMessage) UserMessage() string {
	return w.message
}

func (w *withUserMessage) Error() string {
	return w.cause.Error()
}

func (w *withUserMessage) Cause() error { return w.cause }

func (w *withUserMessage) Unwrap() error { return w.cause }

func (w *withUserMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "user_message=%q %+v\n", w.message, w.Cause())
			return
		}

		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}

// GetUserMessage returns the first non-empty client message in err's chain, or
// "".
func GetUserMessage(err error) string {
	var msg string

	doUnwrap(err, func(e error) bool {
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
