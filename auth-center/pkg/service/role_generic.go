package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/auth-center/api"
	"github.com/runtime-radar/runtime-radar/auth-center/pkg/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type RoleGeneric struct {
	api.UnimplementedRoleControllerServer
	RoleRepository database.RoleRepository
}

func (rg *RoleGeneric) Read(ctx context.Context, roleReq *api.ReadRoleReq) (resp *api.ReadRoleResp, err error) {
	id, err := uuid.Parse(roleReq.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
	}

	role, err := rg.RoleRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "role does not exist")
		}
		return nil, status.Errorf(codes.Internal, "can't get role: %v", err)
	}

	return fillRoleResp(role), nil
}

func (rg *RoleGeneric) ReadList(ctx context.Context, _ *emptypb.Empty) (resp *api.ReadListRoleResp, err error) {
	resp = &api.ReadListRoleResp{}

	roles, err := rg.RoleRepository.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get roles: %v", err)
	}

	for _, role := range roles {
		resp.Roles = append(
			resp.Roles,
			fillRoleResp(role))
	}

	return resp, err
}
