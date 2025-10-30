package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/lib/errcommon"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccessTokenAuth struct {
	AccessToken
	SecretKey []byte
	Verifier  jwt.Verifier
}

func (a *AccessTokenAuth) Create(ctx context.Context, req *model.CreateAccessTokenReq) (id uuid.UUID, token string, err error) {
	if err := a.Verifier.VerifyPermission(ctx, jwt.PermissionPublicAccessTokens, jwt.ActionCreate); err != nil {
		return uuid.Nil, "", errcommon.PermissionErrorToStatus(err)
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, "", status.Errorf(codes.Unauthenticated, "can't get user id: %v", err)
	}

	if userID != req.UserID {
		msg := "current user != target user"
		return uuid.Nil, "", errcommon.StatusWithReason(codes.PermissionDenied, errcommon.PermissionDenied, msg).Err()
	}

	return a.AccessToken.Create(ctx, req)
}

func (a *AccessTokenAuth) ListPage(ctx context.Context, pageNum, pageSize int, order string) ([]*model.AccessTokenResp, int, error) {
	if err := a.Verifier.VerifyPermission(ctx, jwt.PermissionPublicAccessTokens, jwt.ActionRead); err != nil {
		return nil, 0, errcommon.PermissionErrorToStatus(err)
	}

	return a.AccessToken.ListPage(ctx, pageNum, pageSize, order)
}

func (a *AccessTokenAuth) Delete(ctx context.Context, id uuid.UUID) error {
	if err := a.Verifier.VerifyPermission(ctx, jwt.PermissionPublicAccessTokens, jwt.ActionDelete); err != nil {
		return errcommon.PermissionErrorToStatus(err)
	}

	at, err := a.AccessToken.GetByID(ctx, id)
	if err != nil {
		return err
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "can't get user id: %v", err)
	}

	if at.UserID != userID {
		return status.Error(codes.PermissionDenied, "token is not accessible")
	}

	return a.AccessToken.Delete(ctx, id)
}

func (a *AccessTokenAuth) InvalidateAll(ctx context.Context) error {
	if err := a.Verifier.VerifyPermission(ctx, jwt.PermissionInvalidatePublicAccessTokens, jwt.ActionExecute); err != nil {
		return errcommon.PermissionErrorToStatus(err)
	}

	return a.AccessToken.InvalidateAll(ctx)
}

func (a *AccessTokenAuth) GetByID(ctx context.Context, id uuid.UUID) (*model.AccessTokenResp, error) {
	if err := a.Verifier.VerifyPermission(ctx, jwt.PermissionPublicAccessTokens, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}

	at, err := a.AccessToken.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "can't get user id: %v", err)
	}

	if at.UserID != userID {
		return nil, status.Error(codes.PermissionDenied, "token is not accessible")
	}

	return at, nil
}
