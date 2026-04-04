package validate

import (
	stderrors "errors"
	"fmt"
	"strings"
)

// FieldReason matches protoc-gen-validate *ValidationError types (Field +
// Reason).
type FieldReason interface {
	Field() string
	Reason() string
}

// AllErrorsLister matches protoc-gen-validate *MultiError slice types.
type AllErrorsLister interface {
	AllErrors() []error
}

func friendlyValidationMessage(err error) string {
	if err == nil {
		return ""
	}

	var lister AllErrorsLister
	if stderrors.As(err, &lister) {
		errs := lister.AllErrors()

		parts := make([]string, 0, len(errs))
		for _, e := range errs {
			if e == nil {
				continue
			}

			parts = append(parts, formatSingleViolation(e))
		}

		if len(parts) > 0 {
			return strings.Join(parts, "\n")
		}
	}

	return formatSingleViolation(err)
}

func formatSingleViolation(err error) string {
	var fr FieldReason
	if stderrors.As(err, &fr) {
		return fmt.Sprintf("%s: %s", fr.Field(), fr.Reason())
	}

	return err.Error()
}
