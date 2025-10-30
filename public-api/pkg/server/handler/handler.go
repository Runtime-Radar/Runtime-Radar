package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	JSONContentType = "application/json"
)

// marshalOptions returns a protojson.MarshalOptions configuration
// necessary for correctly parsing protojson and maintaining a consistent style
// with other services. It enables the use of proto field names and emits
// unpopulated fields during marshaling.
func marshalOptions() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
}

// SendJSONResp marshals (regular) message given in resp as JSON to w.
func SendJSONResp(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", JSONContentType)
	jsonResp(w, resp)
}

// SendProtoResp marshals (proto) message given in m as JSON to w.
func SendProtoResp(w http.ResponseWriter, m proto.Message) {
	w.Header().Set("Content-Type", JSONContentType)
	protoResp(w, m)
}

// ErrorJSONResp writes a JSON response with status 500 (Internal Server Error).
func ErrorJSONResp(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", JSONContentType)
	w.WriteHeader(http.StatusInternalServerError)
	jsonResp(w, resp)
}

// StatusJSONResp marshals (proto) status error st to w and sets HTTP response code according to status code mapping from statusToHTTP.
func StatusJSONResp(w http.ResponseWriter, st *status.Status) {
	w.Header().Set("Content-Type", JSONContentType)
	w.WriteHeader(statusToHTTP(st.Code()))
	protoResp(w, st.Proto())
}

func jsonResp(w http.ResponseWriter, resp any) {
	log.Debug().Msgf("Sending JSON response: %+v", resp)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Just log, and don't try to override inflight response, because HTTP response code (with other headers)
		// and some data may already be sent, or writer may become unavailable, and in most cases it doesn't make any sense,
		// as it could lead to even more errors and unexpected behavior on client side. Fortunately, this part (almost) never breaks.
		log.Error().Msgf("Can't encode response: %v", err)
	}
}

func protoResp(w http.ResponseWriter, m proto.Message) {
	log.Debug().Msgf("Sending proto message as JSON: %+v", m)
	mo := marshalOptions()
	b, err := mo.Marshal(m)
	if err != nil {
		// Just log, and don't try to override inflight response, because HTTP response code (with other headers)
		// and some data may already be sent, or writer may become unavailable, and in most cases it doesn't make any sense,
		// as it could lead to even more errors and unexpected behavior on client side. Fortunately, this part (almost) never breaks.
		log.Error().Msgf("Can't encode response: %v", err)
	}

	if _, err := w.Write(b); err != nil {
		// Just log, and don't try to override inflight response, because HTTP response code (with other headers)
		// and some data may already be sent, or writer may become unavailable, and in most cases it doesn't make any sense,
		// as it could lead to even more errors and unexpected behavior on client side. Fortunately, this part (almost) never breaks.
		log.Error().Msgf("Can't send response: %v", err)
	}
}

// statusToHTTP converts a gRPC error code into the corresponding HTTP response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func statusToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		log.Warn().Msgf("Unknown gRPC error code: %v", code)
		return http.StatusInternalServerError
	}
}
