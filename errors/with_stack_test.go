package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestWithStack_nilErr(t *testing.T) {
	assert.Nil(t, errors.WithStack(nil))
}

func TestWithStack(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantUnwrap   error
		wantStackLen int
	}

	root := errors.New("root")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "adds-stack",
			args:    Args{err: root},
			expects: Expects{wantUnwrap: root, wantStackLen: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithStack(tc.args.err)
			require.NotNil(t, got)
			assert.Equal(t, root.Error(), got.Error())
			assert.Equal(t, tc.expects.wantUnwrap, errors.Unwrap(got))

			var stackErr errors.StackError
			require.True(t, errors.As(got, &stackErr))
			assert.GreaterOrEqual(t, len(stackErr.StackTrace()), tc.expects.wantStackLen)
			assert.Contains(t, fmt.Sprintf("%+v", got), root.Error())
		})
	}
}

func TestWithStack_Error(t *testing.T) {
	type Args struct {
		inner error
	}
	type Expects struct {
		want string
	}

	root := errors.New("root")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "delegates-to-fundamental",
			args:    Args{inner: root},
			expects: Expects{want: "root"},
		},
		{
			name:    "delegates-to-plain-err",
			args:    Args{inner: &errNoUnwrap{msg: "db timeout"}},
			expects: Expects{want: "db timeout"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithStack(tc.args.inner)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, tc.expects.want, tc.args.inner.Error())
		})
	}
}

func TestWithStack_Format(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS         string
		wantQ         string
		wantV         string
		vPlusContains []string
	}

	inner := &errNoUnwrap{msg: "db timeout"}
	err := errors.WithStack(inner)

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "with-stack",
			args: Args{err: err},
			expects: Expects{
				wantS:         "db timeout",
				wantQ:         `"db timeout"`,
				wantV:         "db timeout",
				vPlusContains: []string{"db timeout"},
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
