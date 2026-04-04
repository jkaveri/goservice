package errors_test

import (
	"strings"
	"testing"

	errors "github.com/jkaveri/goservice/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Join(t *testing.T) {
	assert.Nil(t, errors.Join())

	errs := []error{
		errors.New("err 1"),
		errors.New("err 2"),
		nil,
		errors.New("err 3"),
		errors.New("err 4"),
	}
	arr := make([]string, 0)

	for i := range errs {
		if errs[i] == nil {
			continue
		}

		arr = append(arr, errs[i].Error())
	}

	expected := strings.Join(arr, "\n")

	err := errors.Join(errs...)

	assert.Equal(t, expected, err.Error())
}

func Test_Join_Is(t *testing.T) {
	errs := []error{
		errors.New("err 1"),
		errors.New("err 2"),
		errors.New("err 3"),
		errors.New("err 4"),
	}

	err := errors.Join(errs...)

	for i := range errs {
		assert.Equal(t, true, errors.Is(err, errs[i]), "err should equal: %s", errs[i].Error())
	}
}
