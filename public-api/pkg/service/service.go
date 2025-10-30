package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
)

const (
	DefaultPageNum  = 1
	DefaultPageSize = 10
	DefaultOrder    = "created_at desc"
)

// gRPC errdetails.ErrorInfo.Reason codes used in service responses.
const (
	NameMustBeUnique = "NAME_MUST_BE_UNIQUE"
)

func userIDFromContext(ctx context.Context) (uuid.UUID, error) {
	token, err := jwt.UnverifiedTokenFromContext(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(token.GetUserID())
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't parse user id UUID: %w", err)
	}

	return userID, nil
}
