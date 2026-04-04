package validate

import (
	lib "github.com/go-playground/validator/v10"
	"github.com/jkaveri/goservice/errorcode"
	errors "github.com/jkaveri/goservice/errors"
)

func ValidateStruct(val interface{}) error {
	err := validatorInstance.Struct(val)
	if v, ok := err.(lib.ValidationErrors); ok {
		return errorcode.InvalidRequest(
			validatorInstance.Translate(v),
		)
	}

	if v, ok := err.(*lib.InvalidValidationError); ok {
		return v
	}

	if err != nil {
		return errors.Wrap(err, "failed to validate struct")
	}

	return nil
}
