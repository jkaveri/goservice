package errors_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func Test_WithExtra(t *testing.T) {
	err := errors.New("test")

	errWithMetadata := errors.WithMetadata(err, map[string]interface{}{
		"Email": "john@email.com",
	})

	metadataGetter, ok := errWithMetadata.(errors.MetadataError)
	require.True(t, ok)

	assert.True(t, ok)
	assert.Equal(t, map[string]interface{}{
		"Email": "john@email.com",
	}, metadataGetter.Metadata())

	assert.Equal(t, err, errors.Unwrap(errWithMetadata))
	assert.Equal(t, "test", errWithMetadata.Error())
	verbose := fmt.Sprintf("%+v", errWithMetadata)
	assert.Contains(t, verbose, "test")
	assert.Contains(t, verbose, "metadata:")
	assert.Contains(t, verbose, "Email:")
	assert.Contains(t, verbose, "john@email.com")
	assert.Contains(t, fmt.Sprintf("%s", errWithMetadata), "test")
	assert.Contains(t, fmt.Sprintf("%q", errWithMetadata), "test")
}

func Test_GetMetadata(t *testing.T) {
	assert.Nil(t,
		errors.Metadata(
			errors.New("test"),
		),
	)

	assert.Equal(t,
		map[string]any{
			"Email": "john@email.com",
		},
		errors.Metadata(
			errors.WithMetadata(
				errors.New("test"),
				map[string]any{
					"Email": "john@email.com",
				},
			),
		),
	)

	assert.Equal(t,
		map[string]any{
			"Email": "john@email.com",
		},
		errors.Metadata(
			errors.Wrap(
				errors.WithMetadata(
					errors.New("test"),
					map[string]any{
						"Email": "john@email.com",
					},
				),
				"wrapped",
			),
		),
	)
}

func TestWithMetadata_Error(t *testing.T) {
	type Args struct {
		inner error
	}
	type Expects struct {
		want string
	}

	inner := errors.New("db down")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "delegates-to-inner",
			args:    Args{inner: inner},
			expects: Expects{want: "db down"},
		},
		{
			name: "delegates-through-message-wrapper",
			args: Args{
				inner: errors.WithMessage(errors.New("leaf"), "annot"),
			},
			expects: Expects{want: "annot: leaf"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithMetadata(tc.args.inner, map[string]any{"key": "val"})
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestWithMetadata_Format(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS           string
		wantQ           string
		wantV           string
		wantVPlus       string
		wantVPlusPrefix string
		vPlusContains   []string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "leaf-with-metadata",
			args: Args{
				err: errors.WithMetadata(
					errors.New("inner"),
					map[string]any{"key": "val"},
				),
			},
			expects: Expects{
				wantS:     "inner",
				wantQ:     "\"inner\"",
				wantV:     "inner",
				wantVPlus: "inner\nmetadata:\tkey: val",
			},
		},
		{
			name: "empty-metadata-map-omits-block",
			args: Args{
				err: errors.WithMetadata(errors.New("only"), map[string]any{}),
			},
			expects: Expects{
				wantS:     "only",
				wantQ:     "\"only\"",
				wantV:     "only",
				wantVPlus: "only",
			},
		},
		{
			name: "duplicate-key-style-with-sorted-keys",
			args: Args{
				err: errors.WithMetadata(
					errors.WithMessage(
						errors.New(`duplicate key value violates unique constraint "tenants_slug_key"`),
						"slug already in use",
					),
					map[string]any{
						"tenant_id": "t_7a2",
						"slug":      "acme",
					},
				),
			},
			expects: Expects{
				wantS: `slug already in use: duplicate key value violates unique constraint "tenants_slug_key"`,
				wantQ: `"slug already in use: duplicate key value violates unique constraint \"tenants_slug_key\""`,
				wantV: `slug already in use: duplicate key value violates unique constraint "tenants_slug_key"`,
				wantVPlusPrefix: "slug already in use\nduplicate key value violates unique constraint \"tenants_slug_key\"",
				vPlusContains: []string{
					"metadata:",
					"slug: acme",
					"tenant_id: t_7a2",
					"errors.WithMessage",
					"runtime.",
				},
			},
		},
		{
			name: "nested-metadata-each-layer",
			args: Args{
				err: errors.WithMetadata(
					errors.WithMetadata(
						errors.New("db timeout"),
						map[string]any{"shard": "0"},
					),
					map[string]any{"trace_id": "tr_abc"},
				),
			},
			expects: Expects{
				wantS:     "db timeout",
				wantQ:     "\"db timeout\"",
				wantV:     "db timeout",
				wantVPlus: "db timeout\nmetadata:\tshard: 0\nmetadata:\ttrace_id: tr_abc",
			},
		},
		{
			name: "metadata-on-grpc-style-handler-wrap",
			args: Args{
				err: errors.WithMetadata(
					errors.Wrap(errors.New("connection refused"), "billing.Charge"),
					map[string]any{
						"grpc_code": "UNAVAILABLE",
						"region":    "us-east-1",
					},
				),
			},
			expects: Expects{
				wantS:           "billing.Charge: connection refused",
				wantQ:           "\"billing.Charge: connection refused\"",
				wantV:           "billing.Charge: connection refused",
				wantVPlusPrefix: "billing.Charge\n\tconnection refused",
				vPlusContains: []string{
					"metadata:",
					"grpc_code: UNAVAILABLE",
					"region: us-east-1",
					"errors.Wrap",
					"runtime.",
				},
			},
		},
		{
			name: "background-job-context-deadline",
			args: Args{
				err: errors.WithMetadata(
					fmt.Errorf("run backup: %w", context.DeadlineExceeded),
					map[string]any{
						"job_id":  "bk_9",
						"attempt": 2,
					},
				),
			},
			expects: Expects{
				wantS:     "run backup: context deadline exceeded",
				wantQ:     "\"run backup: context deadline exceeded\"",
				wantV:     "run backup: context deadline exceeded",
				wantVPlus: "run backup: context deadline exceeded\nmetadata:\tattempt: 2\tjob_id: bk_9",
			},
		},
		{
			name: "auth-gateway-with-message-chain",
			args: Args{
				err: errors.WithMetadata(
					errors.WithMessage(
						errors.New("token is malformed"),
						"authorize request",
					),
					map[string]any{
						"kid":  "sig_live_01",
						"path": "/v1/graphql",
					},
				),
			},
			expects: Expects{
				wantS:           "authorize request: token is malformed",
				wantQ:           "\"authorize request: token is malformed\"",
				wantV:           "authorize request: token is malformed",
				wantVPlusPrefix: "authorize request\ntoken is malformed",
				vPlusContains: []string{
					"metadata:",
					"kid: sig_live_01",
					"path: /v1/graphql",
					"errors.WithMessage",
					"runtime.",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.wantS, fmt.Sprintf("%s", tc.args.err))
			assert.Equal(t, tc.expects.wantQ, fmt.Sprintf("%q", tc.args.err))
			assert.Equal(t, tc.expects.wantV, fmt.Sprintf("%v", tc.args.err))

			gotPlus := fmt.Sprintf("%+v", tc.args.err)
			switch {
			case tc.expects.wantVPlusPrefix != "":
				assert.True(t, len(gotPlus) >= len(tc.expects.wantVPlusPrefix))
				assert.Equal(t, tc.expects.wantVPlusPrefix, gotPlus[:len(tc.expects.wantVPlusPrefix)])
				for _, sub := range tc.expects.vPlusContains {
					assert.Contains(t, gotPlus, sub)
				}
			default:
				assert.Equal(t, tc.expects.wantVPlus, gotPlus)
			}
		})
	}
}
