package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/validate"
)

func TestValidator_Translate(t *testing.T) {
	type Company struct {
		Name string `validate:"required"`
	}

	type User struct {
		UserID   string `validate:"shortuuid"`
		Username string `validate:"required_without=UserID"`
		Company  Company
	}

	validator := validate.NewValidator()

	err := validator.Struct(&User{
		UserID:   "idk",
		Username: "",
		Company:  Company{},
	})

	translated := validator.Translate(err)

	assert.Contains(t, translated, "User.UserID: UserID is not valid short uuid")
	assert.Contains(t, translated, "User.Company.Name: Name is a required field")

	err = validator.Struct(&User{
		UserID:   "",
		Username: "",
		Company: Company{
			Name: "henry",
		},
	})

	translated = validator.Translate(err)
	assert.Contains(t, translated, "User.Username: Username and UserID are both empty")

	assert.PanicsWithValue(t,
		"shortuuid validation can only be applied to string",
		func() {
			err = validator.Struct(&struct {
				Age int `validate:"shortuuid"`
			}{
				Age: 10,
			})
		},
	)
}
