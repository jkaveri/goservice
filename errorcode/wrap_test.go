package errorcode_test

import (
	stderrors "errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/errorcode"
	svcerrors "github.com/jkaveri/goservice/errors"
)

func TestWrap_nil_error_returns_nil(t *testing.T) {
	got := errorcode.Wrap(nil, errorcode.CodeNotFound, "annotation")
	assert.Nil(t, got)
}

func TestWrapf_nil_error_returns_nil(t *testing.T) {
	got := errorcode.Wrapf(nil, errorcode.CodeNotFound, "annotation %d", 1)
	assert.Nil(t, got)
}

func TestWrap_tags_code_and_preserves_errors_is(t *testing.T) {
	type Args struct {
		base error
		code string
		msg  string
	}
	type Expects struct {
		wantCode string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "eof-with-not-found-code",
			args: Args{
				base: io.EOF,
				code: errorcode.CodeNotFound,
				msg:  "resource missing",
			},
			expects: Expects{wantCode: errorcode.CodeNotFound},
		},
		{
			name: "coded-inner-still-uses-outer-code",
			args: Args{
				base: errorcode.NotFound("widget absent"),
				code: errorcode.CodeInvalidRequest,
				msg:  "bad query",
			},
			expects: Expects{wantCode: errorcode.CodeInvalidRequest},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errorcode.Wrap(tc.args.base, tc.args.code, tc.args.msg)
			assert.NotNil(t, got)
			assert.Equal(t, tc.expects.wantCode, svcerrors.Code(got))
			assert.True(t, stderrors.Is(got, tc.args.base), "expected errors.Is chain to include base: %v", got)
			assert.Contains(t, got.Error(), tc.args.msg)
		})
	}
}

func TestWrapf_tags_code_formats_message_and_preserves_errors_is(t *testing.T) {
	type Args struct {
		base   error
		code   string
		format string
		args   []any
	}
	type Expects struct {
		wantCode       string
		wantMsgSubstr string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "eof-with-formatted-annotation",
			args: Args{
				base:   io.EOF,
				code:   errorcode.CodeTimeout,
				format: "deadline exceeded after %dms",
				args:   []any{500},
			},
			expects: Expects{
				wantCode:       errorcode.CodeTimeout,
				wantMsgSubstr: "deadline exceeded after 500ms",
			},
		},
		{
			name: "literal-format-without-args",
			args: Args{
				base:   io.EOF,
				code:   errorcode.CodeUnavailable,
				format: "upstream offline",
				args:   nil,
			},
			expects: Expects{
				wantCode:       errorcode.CodeUnavailable,
				wantMsgSubstr: "upstream offline",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errorcode.Wrapf(tc.args.base, tc.args.code, tc.args.format, tc.args.args...)
			assert.NotNil(t, got)
			assert.Equal(t, tc.expects.wantCode, svcerrors.Code(got))
			assert.True(t, stderrors.Is(got, tc.args.base), "expected errors.Is chain to include base: %v", got)
			assert.Contains(t, got.Error(), tc.expects.wantMsgSubstr)
		})
	}
}
