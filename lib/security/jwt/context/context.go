package context

import (
	"context"
	"errors"
	"fmt"

	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
)

type contextKey string

const userIDContextKey contextKey = "userIDFromToken"

func WithEmptyUserID(ctx context.Context) context.Context {
	userID := ""

	return context.WithValue(ctx, userIDContextKey, &userID)
}

func SetUserID(ctx context.Context) error {
	userID, ok := ctx.Value(userIDContextKey).(*string)
	if !ok || userID == nil {
		return errors.New("userID is not set in context")
	}

	token, err := jwt.UnverifiedTokenFromContext(ctx)
	if err != nil {
		*userID = ""
		return fmt.Errorf("can't parse token to get userID: %w", err)
	}

	*userID = token.GetUserID()

	return nil
}

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDContextKey).(*string)
	if !ok || userID == nil {
		return "", errors.New("userID is not set in context")
	}

	return *userID, nil
}
