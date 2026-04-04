package check_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/check"
)

func TestPanicIfError(t *testing.T) {
	assert.Panics(t, func() {
		check.PanicIfError(errors.New("test"))
	})
}

func TestPanicIfNil(t *testing.T) {
	assert.Panics(t, func() {
		check.PanicIfNil(nil, "panic because nil")
	})
}

func TestPanicIf(t *testing.T) {
	assert.Panics(t, func() {
		check.PanicIf(true, "panic because true")
	})

	assert.NotPanics(t, func() {
		check.PanicIf(false, "panic because false")
	})
}

func TestTenary(t *testing.T) {
	assert.Equal(t, 11, check.Ternary(11 > 0, 11, 0))
	assert.Equal(t, 12, check.Ternary(11 > 12, 11, 12))
}
