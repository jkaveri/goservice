package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	errors "github.com/jkaveri/goservice/errors"
)

type errNoUnwrap struct {
	msg string
}

func (e *errNoUnwrap) Error() string { return e.msg }

func TestWalkErrorChain(t *testing.T) {
	type Args struct {
		err error
	}

	type Expects struct {
		wantMessages []string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "nil-skip-callback",
			args:    Args{err: nil},
			expects: Expects{wantMessages: nil},
		},
		{
			name:    "leaf-without-unwrap",
			args:    Args{err: errors.New("leaf")},
			expects: Expects{wantMessages: []string{"leaf"}},
		},
		{
			name: "unwrap-chain",
			args: Args{
				err: errors.Wrap(errors.New("root"), "outer"),
			},
			expects: Expects{
				wantMessages: []string{"outer: root", "root", "root"},
			},
		},
		{
			name: "stops-without-unwrap",
			args: Args{
				err: &errNoUnwrap{msg: "solo"},
			},
			expects: Expects{wantMessages: []string{"solo"}},
		},
		{
			name: "join-two-leaves",
			args: Args{
				err: errors.Join(errors.New("a"), errors.New("b")),
			},
			expects: Expects{
				wantMessages: []string{"a\nb", "a", "b"},
			},
		},
		{
			name: "join-skips-nil",
			args: Args{
				err: errors.Join(errors.New("a"), nil, errors.New("b")),
			},
			expects: Expects{
				wantMessages: []string{"a\nb", "a", "b"},
			},
		},
		{
			name: "join-all-nil",
			args: Args{
				err: errors.Join(nil, nil),
			},
			expects: Expects{wantMessages: nil},
		},
		{
			name: "join-walks-each-branch",
			args: Args{
				err: errors.Join(
					errors.Wrap(errors.New("inner"), "outer"),
					errors.New("sibling"),
				),
			},
			expects: Expects{
				wantMessages: []string{
					"outer: inner\nsibling",
					"outer: inner",
					"inner",
					"inner",
					"sibling",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got []string

			errors.WalkErrorChain(tc.args.err, func(e error) bool {
				got = append(got, e.Error())

				return false
			})

			assert.Equal(t, tc.expects.wantMessages, got)
		})
	}
}

func TestWalkErrorChainStopOnTrue(t *testing.T) {
	type Args struct {
		err error
	}

	type Expects struct {
		wantMessages []string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "single-unwrap-chain",
			args: Args{
				err: errors.Wrap(errors.New("inner"), "outer"),
			},
			expects: Expects{
				wantMessages: []string{"outer: inner", "inner"},
			},
		},
		{
			name: "join-stops-before-second-leaf",
			args: Args{
				err: errors.Join(errors.New("a"), errors.New("b")),
			},
			expects: Expects{wantMessages: []string{"a\nb", "a"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got []string

			errors.WalkErrorChain(tc.args.err, func(e error) bool {
				got = append(got, e.Error())

				return len(got) == 2
			})

			assert.Equal(t, tc.expects.wantMessages, got)
		})
	}
}
