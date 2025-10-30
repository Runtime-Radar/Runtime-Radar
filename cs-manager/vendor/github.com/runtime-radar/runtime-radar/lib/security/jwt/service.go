package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// 87600 hours == 10 years.
// This time was chosen as a long enough period for "infinite" token.
const tenYears = 87600 * time.Hour

func GenerateServiceToken(key []byte, name string, rolePermissions *RolePermissions) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key is empty")
	}

	t := &Token{
		Username: name,
		Role: &Role{
			RoleName:        "service",
			RolePermissions: rolePermissions,
		},
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: &jwtv5.NumericDate{Time: time.Now().Add(tenYears)},
			IssuedAt:  &jwtv5.NumericDate{Time: time.Now()},
		},
	}

	return SignHS256(t, key)
}
