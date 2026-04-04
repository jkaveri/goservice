package validate_test

import (
	"testing"

	"github.com/jkaveri/goservice/errorcode"
	"github.com/jkaveri/goservice/validate"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateStruct(t *testing.T) {
	type User struct {
		UserID   string `validate:"shortuuid"`
		Username string `validate:"required"`
	}

	assert.NoError(t, validate.ValidateStruct(&User{Username: "test"}))
	assert.True(
		t,
		errorcode.IsInvalidRequest(validate.ValidateStruct(&User{Username: ""})),
	)
	assert.True(
		t,
		errorcode.IsInvalidRequest(validate.ValidateStruct(&User{Username: "test", UserID: "invalid"})),
	)
}
