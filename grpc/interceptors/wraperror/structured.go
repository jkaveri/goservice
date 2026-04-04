package wraperror

import (
	"encoding/json"
	"net/http"

	errors "github.com/jkaveri/goservice/errors"
	"github.com/jkaveri/ramda"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/jkaveri/goservice/errorcode"
)

const MetadataKeyAppErrorCode = "@app_error_code"

func ToStructured(err error) *StructuredError {
	if err == nil {
		return nil
	}

	if in, ok := err.(*StructuredError); ok {
		return in
	}

	code := getErrorCode(err)
	metadata := errors.Metadata(err)

	msg := errors.GetUserMessage(err)
	if msg == "" {
		msg = err.Error()
	}

	// fallback for grpc-gateway handler
	// see grpcx/gateway/error_handler.go
	metadata[MetadataKeyAppErrorCode] = code

	return &StructuredError{
		Code:         code,
		ErrorMessage: msg,
		Metadata:     metadata,
		err:          err,
	}
}

type (
	StructuredErrorAlias StructuredError
	StructuredError      struct {
		Code string `json:"code"`
		// fallback for http client error serialize
		ErrorMessage string                 `json:"error"`
		Metadata     map[string]interface{} `json:"metadata"`

		err error
	}
)

func (e *StructuredError) Error() string {
	if e.err == nil {
		return e.ErrorMessage
	}

	return e.err.Error()
}

// StatusCode http status code
func (e *StructuredError) StatusCode() int {
	return CodeToHTTPStatus(e.Code)
}

// GRPCStatus grpc status
func (e *StructuredError) GRPCStatus() *status.Status {
	s := status.New(
		CodeToGRPC(e.Code),
		e.ErrorMessage,
	)

	if len(e.Metadata) > 0 {
		details, _ := structpb.NewValue(e.Metadata)

		s, _ = s.WithDetails(details)
	}

	return s
}

// MarshalJSON marshal error to json
// this implements json marshaler to support error handling of truss
func (e *StructuredError) MarshalJSON() ([]byte, error) {
	return json.Marshal((*StructuredErrorAlias)(e))
}

// UnmarshalJSON ...
func (e *StructuredError) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, (*StructuredErrorAlias)(e))
}

func getErrorCode(err error) string {
	return ramda.Default(
		errors.Code(err),
		errorcode.CodeInternalServer,
	)
}

func CodeToHTTPStatus(code string) int {
	switch code {
	case errorcode.CodeNone:
		return http.StatusOK
	case errorcode.CodeInvalidRequest:
		return http.StatusBadRequest
	case errorcode.CodeNotFound:
		return http.StatusNotFound
	case errorcode.CodeUnauthorized:
		return http.StatusForbidden
	case errorcode.CodeNotAuthenticated:
		return http.StatusUnauthorized
	case errorcode.CodeDuplicated:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func CodeFromHTTPStatus(status int) string {
	switch status {
	case http.StatusOK:
		return errorcode.CodeNone
	case http.StatusBadRequest:
		return errorcode.CodeInvalidRequest
	case http.StatusNotFound:
		return errorcode.CodeNotFound
	case http.StatusForbidden:
		return errorcode.CodeUnauthorized
	case http.StatusUnauthorized:
		return errorcode.CodeNotAuthenticated
	case http.StatusConflict:
		return errorcode.CodeDuplicated
	default:
		return errorcode.CodeInternalServer
	}
}

func CodeToGRPC(code string) codes.Code {
	switch code {
	case errorcode.CodeNone:
		return codes.OK
	case errorcode.CodeInvalidRequest:
		return codes.InvalidArgument
	case errorcode.CodeNotFound:
		return codes.NotFound
	case errorcode.CodeUnauthorized:
		return codes.PermissionDenied
	case errorcode.CodeNotAuthenticated:
		return codes.Unauthenticated
	case errorcode.CodeDuplicated:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}

func CodeFromGRPC(code codes.Code) string {
	switch code {
	case codes.OK:
		return errorcode.CodeNone
	case codes.InvalidArgument:
		return errorcode.CodeInvalidRequest
	case codes.NotFound:
		return errorcode.CodeNotFound
	case codes.PermissionDenied:
		return errorcode.CodeUnauthorized
	case codes.Unauthenticated:
		return errorcode.CodeNotAuthenticated
	case codes.AlreadyExists:
		return errorcode.CodeDuplicated
	default:
		return errorcode.CodeInternalServer
	}
}

func GenericMessageFromCode(err error) string {
	code := errors.Code(err)

	switch code {
	case errorcode.CodeNone:
		return "ok"
	case errorcode.CodeInvalidRequest:
		return "invalid request"
	case errorcode.CodeNotFound:
		return "not found"
	case errorcode.CodeUnauthorized:
		return "forbidden"
	case errorcode.CodeNotAuthenticated:
		return "not authenticated"
	case errorcode.CodeDuplicated:
		return "conflict"
	default:
		return "internal server error"
	}
}
