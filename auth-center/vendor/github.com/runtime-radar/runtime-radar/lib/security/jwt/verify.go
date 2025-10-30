package jwt

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/runtime-radar/runtime-radar/lib/security/cipher"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrTokenExpired     = jwt.ErrTokenExpired
	ErrSignatureInvalid = jwt.ErrSignatureInvalid
)

type Verifier interface {
	VerifyPermission(ctx context.Context, p PermissionType, as ...Action) error
}

type KeyVerifier struct {
	key []byte
}

func NewKeyVerifier(h string) (*KeyVerifier, []byte, error) {
	key, err := cipher.ParseKey(h)
	if err != nil {
		return nil, nil, err
	}
	return &KeyVerifier{key}, key, nil
}

// VerifyPermission is a basic permissions verifier for CS token.
// All internal errors are returned as ErrorUnauthenticated to be sure
// that there are no internal implementation details are leaked to response.
func (v *KeyVerifier) VerifyPermission(ctx context.Context, pt PermissionType, as ...Action) error {
	if len(as) == 0 {
		panic("empty Actions is not allowed here")
	}
	t, err := TokenFromContext(ctx, v.key)
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return ErrTokenExpired

	case errors.Is(err, jwt.ErrSignatureInvalid):
		return ErrSignatureInvalid

	case err != nil:
		return fmt.Errorf("%w: %w", ErrUnauthenticated, err)
	}

	if t.Role == nil {
		return ErrUnauthenticated
	}

	if !t.Role.RolePermissions.GetPermission(pt).containsAll(as...) {
		return ErrPermissionDenied
	}

	return nil
}
