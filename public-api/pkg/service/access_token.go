package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/model"
)

type AccessToken interface {
	Create(ctx context.Context, req *model.CreateAccessTokenReq) (id uuid.UUID, token string, err error)
	ListPage(ctx context.Context, pageNum, pageSize int, order string) (tokens []*model.AccessTokenResp, total int, err error)
	Delete(ctx context.Context, id uuid.UUID) error
	InvalidateAll(ctx context.Context) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.AccessTokenResp, error)
}
