package errcommon

import (
	"errors"
	"fmt"

	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gRPC errdetails.ErrorInfo.Reason codes used in service responses.
const (
	Unauthenticated    = "UNAUTHENTICATED"
	PermissionDenied   = "PERMISSION_DENIED"
	IncorrectSignature = "INCORRECT_SIGNATURE"
	TokenExpired       = "ACCESS_TOKEN_EXPIRED"
	Internal           = "INTERNAL"
	CodeBadRequest     = "BAD_REQUEST"
	CodeNotFound       = "NOT_FOUND"
)

// StatusWithReason returns grpc status with reason specified
func StatusWithReason(code codes.Code, reason, msg string) *status.Status {
	st := status.New(code, msg)
	st, err := st.WithDetails(&errdetails.ErrorInfo{Reason: reason})
	if err != nil { // normally should not happen
		panic(fmt.Errorf("unexpected error attaching details to status: %w", err))
	}

	return st
}

// ReasonFromStatus returns grpc error's reason if it's provided
func ReasonFromStatus(st *status.Status) (string, bool) {
	for _, d := range st.Details() {
		if errinfo, ok := d.(*errdetails.ErrorInfo); ok {
			return errinfo.Reason, true
		}
	}

	return "", false
}

// PermissionErrorToStatus returns grpc error status with reason from error
func PermissionErrorToStatus(err error) error {
	switch {
	case errors.Is(err, jwt.ErrUnauthenticated):
		return StatusWithReason(codes.Unauthenticated, Unauthenticated, "request is not authenticated").Err()

	case errors.Is(err, jwt.ErrSignatureInvalid):
		return StatusWithReason(codes.Unauthenticated, IncorrectSignature, "signature is incorrect").Err()

	case errors.Is(err, jwt.ErrTokenExpired):
		return StatusWithReason(codes.Unauthenticated, TokenExpired, "token is expired").Err()

	case errors.Is(err, jwt.ErrPermissionDenied):
		return StatusWithReason(codes.PermissionDenied, PermissionDenied, "request is not authorized").Err()
	}

	return StatusWithReason(codes.Internal, Internal, "internal error during token parsing").Err()
}
