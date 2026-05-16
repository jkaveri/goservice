package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	errors "github.com/jkaveri/goservice/errors"
)

func Test_WithCode(t *testing.T) {
	err := errors.New("test")

	errWithCode := errors.WithCode(err, "T123")

	assert.Nil(t, errors.WithCode(nil, "T123"))

	codeContainer, ok := errWithCode.(errors.CodeError)
	assert.True(t, ok)
	assert.Equal(t, "T123", codeContainer.Code())
	assert.True(t, errors.Is(errWithCode, err))
	assert.Equal(t, "[T123] test", errWithCode.Error())
	assert.Contains(t, fmt.Sprintf("%+v", errWithCode), "[T123]")
	assert.Contains(t, fmt.Sprintf("%+v", errWithCode), "test")
}

func TestWithCode_Error(t *testing.T) {
	type Args struct {
		err  error
		code string
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
			name: "delegates-to-simple-err",
			args: Args{
				err:  errors.New("boom"),
				code: "E1",
			},
			expects: Expects{want: "[E1] boom"},
		},
		{
			name: "delegates-to-wrapped-err",
			args: Args{
				err:  errors.Wrap(errors.New("inner"), "outer"),
				code: "E2",
			},
			expects: Expects{want: "[E2] outer: inner"},
		},
		{
			name: "empty-err-message",
			args: Args{
				err:  errors.New(""),
				code: "E3",
			},
			expects: Expects{want: "[E3] "},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithCode(tc.args.err, tc.args.code)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func Test_GetCode(t *testing.T) {
	assert.Equal(t,
		"",
		errors.Code(
			errors.New("test"),
		),
	)

	assert.Equal(t,
		"T123",
		errors.Code(
			errors.WithCode(errors.New("test"), "T123"),
		),
	)

	assert.Equal(t,
		"T123",
		errors.Code(
			errors.Wrap(
				errors.WithCode(errors.New("test"), "T123"),
				"wrapped",
			),
		),
	)
}

func Test_ContainsCode(t *testing.T) {
	assert.False(t,
		errors.ContainsCode(
			errors.New("test"),
			"T123",
		),
	)

	assert.True(t,
		errors.ContainsCode(
			errors.WithCode(errors.New("test"), "T123"),
			"T123",
		),
	)

	assert.True(t,
		errors.ContainsCode(
			errors.Wrap(
				errors.WithCode(errors.New("test"), "T123"),
				"wrapped",
			),
			"T123",
		),
	)
}

func TestWithCode_Format(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS         string
		wantQ         string
		wantV         string
		vPlusContains []string
	}

	inner := errors.New("inner")
	err := errors.WithCode(inner, "E1")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "with-code",
			args: Args{err: err},
			expects: Expects{
				wantS:         "[E1] inner",
				wantQ:         "\"[E1] inner\"",
				wantV:         "[E1] inner",
				vPlusContains: []string{"[E1]", "inner"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.wantS, fmt.Sprintf("%s", tc.args.err))
			assert.Equal(t, tc.expects.wantQ, fmt.Sprintf("%q", tc.args.err))
			assert.Equal(t, tc.expects.wantV, fmt.Sprintf("%v", tc.args.err))

			gotPlus := fmt.Sprintf("%+v", tc.args.err)
			for _, sub := range tc.expects.vPlusContains {
				assert.Contains(t, gotPlus, sub)
			}
		})
	}
}
