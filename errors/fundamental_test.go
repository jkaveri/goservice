package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestNew_nilMessage(t *testing.T) {
	got := errors.New("")
	require.NotNil(t, got)
	assert.Equal(t, "", got.Error())
}

func TestNew(t *testing.T) {
	type Args struct {
		message string
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
			name:    "returns-message",
			args:    Args{message: "something failed"},
			expects: Expects{want: "something failed"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.New(tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Contains(t, fmt.Sprintf("%+v", got), tc.expects.want)
		})
	}
}

func TestErrorf(t *testing.T) {
	type Args struct {
		format string
		args   []any
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
			name: "formats-message",
			args: Args{
				format: "code %d: %s",
				args:   []any{42, "timeout"},
			},
			expects: Expects{want: "code 42: timeout"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Errorf(tc.args.format, tc.args.args...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Contains(t, fmt.Sprintf("%+v", got), tc.expects.want)
		})
	}
}

func TestNew_Error(t *testing.T) {
	type Args struct {
		message string
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
			name:    "returns-message",
			args:    Args{message: "something failed"},
			expects: Expects{want: "something failed"},
		},
		{
			name:    "empty-message",
			args:    Args{message: ""},
			expects: Expects{want: ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.New(tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestErrorf_Error(t *testing.T) {
	type Args struct {
		format string
		args   []any
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
			name: "formats-message",
			args: Args{
				format: "code %d: %s",
				args:   []any{42, "timeout"},
			},
			expects: Expects{want: "code 42: timeout"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Errorf(tc.args.format, tc.args.args...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestNew_Format(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS         string
		wantQ         string
		wantV         string
		vPlusContains []string
	}

	err := errors.New("boom")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "new-error",
			args: Args{err: err},
			expects: Expects{
				wantS:         "boom",
				wantQ:         `"boom"`,
				wantV:         "boom",
				vPlusContains: []string{"boom"},
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
