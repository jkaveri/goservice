package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFieldErr struct {
	field  string
	reason string
}

func (e testFieldErr) Field() string  { return e.field }
func (e testFieldErr) Reason() string { return e.reason }
func (e testFieldErr) Error() string  { return "validation" }

type testMultiErr []error

func (m testMultiErr) Error() string      { return "multi" }
func (m testMultiErr) AllErrors() []error { return m }

func Test_friendlyValidationMessage(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		want string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "nil",
			args:    Args{err: nil},
			expects: Expects{want: ""},
		},
		{
			name: "single-field-reason",
			args: Args{
				err: testFieldErr{
					field:  "Message",
					reason: "value length must be at least 1 runes",
				},
			},
			expects: Expects{want: "Message: value length must be at least 1 runes"},
		},
		{
			name: "multi-newline-separated",
			args: Args{
				err: testMultiErr{
					testFieldErr{field: "Name", reason: "value length must be at least 1 runes"},
					testFieldErr{field: "Age", reason: "value must be greater than 0"},
				},
			},
			expects: Expects{
				want: "Name: value length must be at least 1 runes\nAge: value must be greater than 0",
			},
		},
		{
			name:    "plain-error-fallback",
			args:    Args{err: assert.AnError},
			expects: Expects{want: assert.AnError.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := friendlyValidationMessage(tc.args.err)
			assert.Equal(t, tc.expects.want, got)
		})
	}
}
