package tokens

import (
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/runtime-radar/runtime-radar/auth-center/pkg/model"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
)

const (
	AuthorizationKey = "authorization"
	RefreshTokenSalt = "RefreshTokenSalt"
)

type BaseToken struct {
	UserID                string        `json:"user_id"`
	Username              string        `json:"username"`
	TokenType             jwt.TokenType `json:"token_type"`
	LastPasswordChangedAt int64         `json:"last_password_changed_at"`
	AuthType              string        `json:"auth_type"`

	jwtv5.RegisteredClaims
}

type TokenPair struct {
	AccessTokenHash  string
	RefreshTokenHash string
}

func GenerateTokenPair(
	user model.User,
	tokenKey []byte,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) (tp TokenPair, err error) {
	at := NewAccessToken(user, accessTokenTTL)
	ats, err := signHS256(at, tokenKey)
	if err != nil {
		return TokenPair{}, err
	}

	rt := NewRefreshToken(user, refreshTokenTTL)
	rts, err := signHS256(rt, tokenKey)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessTokenHash:  ats,
		RefreshTokenHash: rts,
	}, nil
}

func signHS256(
	token jwtv5.Claims,
	TokenKey []byte,
) (string, error) {
	builder := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, token)
	signed, err := builder.SignedString(TokenKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}
