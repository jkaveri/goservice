package grpcgateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"slices"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/goservice/grpc/interceptors/wraperror"
	spbstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Detail struct {
	Type  string          `json:"@type"`
	Value json.RawMessage `json:"value"`
}

type Original struct {
	Code    int64    `json:"code"`
	Message string   `json:"message"`
	Details []Detail `json:"details"`
}

// ErrorHandlerFunc is a function that handles errors of grpc gateway
//
// this function is a copy of runtime.DefaultHTTPErrorHandler
// with some modification to support structured error
var ErrorHandler = runtime.ErrorHandlerFunc(func(
	ctx context.Context,
	sm *runtime.ServeMux,
	m runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request, err error,
) {
	log := golog.WithContext(ctx)

	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
	}

	s := status.Convert(err)
	pb := s.Proto()

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")

	contentType := m.ContentType(pb)
	w.Header().Set("Content-Type", contentType)

	if s.Code() == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", s.Message())
	}

	buf, merr := transformErrorResponse(pb, m)
	if merr != nil {
		log.WithError(merr).Error(
			"failed to marshal error message",
			golog.String("grpc_status", s.String()),
		)
		w.WriteHeader(http.StatusInternalServerError)

		if _, werr := io.WriteString(
			w,
			`{"code": 13, "message": "failed to marshal error message"}`,
		); werr != nil {
			log.WithError(werr).Error("failed to write error response body")
		}

		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		log.Debug("server metadata not in context")
	}

	handleForwardResponseServerMetadata(w, sm, md, []string{
		"content-type",
	})

	// RFC 7230 https://tools.ietf.org/html/rfc7230#section-4.1.2
	// Unless the request includes a TE header field indicating "trailers"
	// is acceptable, as described in Section 4.3, a server SHOULD NOT
	// generate trailer fields that it believes are necessary for the user
	// agent to receive.
	doForwardTrailers := requestAcceptsTrailers(r)

	if doForwardTrailers {
		handleForwardResponseTrailerHeader(w, md)
		w.Header().Set("Transfer-Encoding", "chunked")
	}

	st := runtime.HTTPStatusFromCode(s.Code())
	if customStatus != nil {
		st = customStatus.HTTPStatus
	}

	w.WriteHeader(st)

	if _, werr := w.Write(buf); werr != nil {
		log.WithError(werr).Error("failed to write response")
	}

	if doForwardTrailers {
		handleForwardResponseTrailer(w, md)
	}
})

func transformErrorResponse(
	pb *spbstatus.Status,
	marshaler runtime.Marshaler,
) ([]byte, error) {
	data, merr := marshaler.Marshal(pb)
	if merr != nil {
		return nil, merr
	}

	var original Original

	if err := marshaler.Unmarshal(data, &original); err != nil {
		return nil, err
	}

	metadata := convertDetailsToMetadata(original.Details, marshaler)
	appErrCode := extractAppErrCode(metadata)

	body := map[string]interface{}{
		"code":     appErrCode,
		"message":  original.Message,
		"metadata": metadata,
	}

	newData, err := marshaler.Marshal(body)
	if err != nil {
		return nil, err
	}

	return newData, nil
}

func convertDetailsToMetadata(
	details []Detail,
	marshaler runtime.Marshaler,
) map[string]interface{} {
	metadata := make(map[string]interface{})

	for _, item := range details {
		if item.Type == "type.googleapis.com/google.protobuf.Value" {
			var m map[string]interface{}
			if err := marshaler.Unmarshal(item.Value, &m); err != nil {
				return nil
			}

			for k, v := range m {
				metadata[k] = v
			}
		}
	}

	return metadata
}

func extractAppErrCode(metadata map[string]interface{}) string {
	if code, ok := metadata[wraperror.MetadataKeyAppErrorCode]; ok {
		delete(metadata, wraperror.MetadataKeyAppErrorCode)

		return code.(string)
	}

	return ""
}

func handleForwardResponseServerMetadata(
	w http.ResponseWriter,
	_ *runtime.ServeMux,
	md runtime.ServerMetadata,
	excludes []string,
) {
	for k, vs := range md.HeaderMD {
		if slices.Contains(excludes, strings.ToLower(k)) {
			continue
		}

		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
}

func requestAcceptsTrailers(req *http.Request) bool {
	te := req.Header.Get("TE")
	return strings.Contains(strings.ToLower(te), "trailers")
}

func handleForwardResponseTrailerHeader(
	w http.ResponseWriter,
	md runtime.ServerMetadata,
) {
	for k := range md.TrailerMD {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf(
			"%s%s", runtime.MetadataTrailerPrefix, k,
		))

		w.Header().Add("Trailer", tKey)
	}
}

func handleForwardResponseTrailer(
	w http.ResponseWriter,
	md runtime.ServerMetadata,
) {
	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}
